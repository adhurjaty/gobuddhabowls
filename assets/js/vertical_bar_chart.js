import * as d3 from 'd3';
import d3Tip from '@lix/d3-tip';
d3.tip = d3Tip;

// code from http://bl.ocks.org/Caged/6476579

// Creates a bar chart for category breakdown. May expand use case for the future
export class VerticalBarChart {
    // pass height of container, data to graph and ID string of the element on the main page (not the partial)
    constructor(height, data, id) {
        this.height = height;
        this.data = JSON.parse(data);
        // assume that html is structured:
        // <div id="<id>">
        //   <div style="width:100%">
        //   </div>
        // </div>
        this.divContainer = d3.select('#' + id)
                              .select('div')
                              .attr('class', 'vertical-bar-chart')
                              .style('height', height);
        this.svg = this.divContainer.append('svg');

        this.redraw();
        self = this;
        window.addEventListener('resize', function() {
            self.redraw();
        });
    }

    redraw() {
        this.svg.selectAll('*').remove();
        var totalWidth = this.divContainer.node().getBoundingClientRect().width;
        var margin = { top: 30, right: 40, bottom: 30, left: 40 },
            width = totalWidth - margin.left - margin.right,
            height = this.height - margin.top - margin.bottom;

        var formatDollar = d3.format("$0");
        var x = d3.scaleBand()
            .rangeRound([0, width], .1);

        var y = d3.scaleLinear()
            .range([height, 0]);

        var tip = d3.tip()
            .attr('class', 'd3-tip')
            .offset([-10, 0])
            .html(function (d) {
                return "<strong>Cost:</strong> <span style='color:red'>" + d.Value + "</span>";
            })

        this.svg.style("width", "100%")
            .style("height", this.height)
            .append("g")
            .attr("transform", "translate(" + margin.left + "," + 0 + ")");

        this.svg.call(tip);

        x.domain(this.data.map(function (d) { return d.Name; }));
        y.domain([0, d3.max(this.data, function (d) { return d.Value; })]);

        // Add the x Axis
        this.svg.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate(" + margin.right + "," + (height + margin.top) + ")")
            .call(d3.axisBottom(x));

        // Add the y Axis
        this.svg.append("g")
            .attr("class", "y axis")
            .attr("transform", "translate(" + margin.right + "," + margin.top + ")")            
            .call(d3.axisLeft(y))
            .append("text")
            .attr("transform", "rotate(-90)")
            .attr("y", 6)
            .attr("dy", ".71em")
            .style("text-anchor", "end")
            .text("Cost");

        this.svg.selectAll(".bar")
            .data(this.data)
            .enter().append("rect")
            .attr("class", "bar")
            .attr("x", function (d) { return x(d.Name); })
            .attr("width", x.bandwidth())
            .attr("y", function (d) { return y(d.Value); })
            .attr("height", function (d) { return height - y(d.Value); })
            .attr("fill", function(d) { return d.Background; })
            .attr("transform", "translate(" + margin.right + "," + margin.top + ")")            
            .on('mouseover', function(d) {
                tip.show(d);
                d3.select(this).style("fill", function() {
                    return d3.rgb(d3.select(this).style("fill")).darker(0.5);
                });
            })
            .on('mouseout', function(d) {
                tip.hide(d);
                d3.select(this).style("fill", d.Background);
            });
    }
}