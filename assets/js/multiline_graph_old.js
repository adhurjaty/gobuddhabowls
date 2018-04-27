import * as d3 from 'd3';
import d3Tip from '@lix/d3-tip';
d3.tip = d3Tip;
import { d3LineWithLegend } from "./d3linewithlegend.js";

export class MultilineGraph {
    constructor(height, data, id) {
        this.height = height;
        // this.data = JSON.parse(data);
        this.data = this.generateData();

        // assume that html is structured:
        // <div id="<id>">
        //   <div style="width:100%">
        //   </div>
        // </div>
        this.divContainer = d3.select('#' + id)
                              .select('div')
                              .attr('class', 'multiline-graph')
                              .style('height', height);
        this.svg = this.divContainer.append('svg');
        this.margin = {top: 30, right: 10, bottom: 50, left: 60};

        this.createGraph();
        self = this;
        $(window).resize(function() {
            self.chart.width(self.getWidth(self.margin))
                .height(self.getHeight(self.margin));
    
            self.svg
                .attr('width', self.getWidth(self.margin))
                .attr('height', self.getHeight(self.margin))
                .call(self.chart);
        });
    }

    log(text) {
        if (console && console.log) console.log(text);
        return text;
    }

    createGraph() {
        this.chart = d3LineWithLegend()
                    .xAxis.label('Time (ms)')
                    .width(this.getWidth(this.margin))
                    .height(this.getHeight(this.margin))
                    .yAxis.label('Voltage (v)');

        this.svg.datum(this.data);

        this.svg
            .attr('width', this.getWidth(this.margin))
            .attr('height', this.getHeight(this.margin))
            .call(this.chart);


        this.chart.dispatch.on('showTooltip', function(e) {
            var offset = this.divContainer.offset(), // { left: 0, top: 0 }
                left = e.pos[0] + offset.left,
                top = e.pos[1] + offset.top,
                formatter = d3.format(".04f");

            var content = '<h3>' + e.series.label + '</h3>' +
                            '<p>' +
                            '<span class="value">[' + e.point[0] + ', ' + formatter(e.point[1]) + ']</span>' +
                            '</p>';

            nvtooltip.show([left, top], content);
        });

        this.chart.dispatch.on('hideTooltip', function(e) {
            nvtooltip.cleanup();
        });
    }
    
    getWidth(margin) {
        var w = $(window).width() - 20;

        return ( (w - margin.left - margin.right - 20) < 0 ) ? margin.left + margin.right + 2 : w;
    }

    getHeight(margin) {
        var h = this.height - 20;

        return ( h - margin.top - margin.bottom - 20 < 0 ) ? 
                    margin.top + margin.bottom + 2 : h;
    }

    //data
    generateData() {
        var sin = [],
            sin2 = [],
            cos = [],
            cos2 = [],
            r1 = Math.random(),
            r2 = Math.random(),
            r3 = Math.random(),
            r4 = Math.random();

        for (var i = 0; i < 100; i++) {
            sin.push([ i, r1 * Math.sin( r2 +  i / (10 * (r4 + .5) ))]);
            cos.push([ i, r2 * Math.cos( r3 + i / (10 * (r3 + .5) ))]);
            sin2.push([ i, r3 * Math.sin( r1 + i / (10 * (r2 + .5) ))]);
            cos2.push([ i, r4 * Math.cos( r4 + i / (10 * (r1 + .5) ))]);
        }

        return [
            {
                data: sin,
                label: "Sine Wave"
            },
            {
                data: cos,
                label: "Cosine Wave"
            },
            {
                data: sin2,
                label: "Sine2 Wave"
            },
            {
                data: cos2,
                label: "Cosine2 Wave"
            }
        ];
    }

}
