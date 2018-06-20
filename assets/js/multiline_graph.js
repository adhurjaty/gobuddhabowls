// code from https://gist.github.com/d3noob/f8b7f107ba25c21971851728520224cb

import * as d3 from 'd3';
import d3Tip from '@lix/d3-tip';
d3.tip = d3Tip;

export class MultilineGraph {
    constructor(height, data, id) {
        this.height = height;
        this.data = data;

        // assume that html is structured:
        // <div id="<id>">
        //   <div style="width:100%">
        //   </div>
        // </div>
        this.divContainer = d3.select('#' + id)
                              .attr('class', 'multiline-graph')
                              .style('height', height);
        this.svg = this.divContainer.append('svg');

        this.redraw();

        window.addEventListener('resize', function(context) {
            return function() {
                context.redraw()
            };
        }(this));
    }

    redraw() {
        var self = this;
        this.svg.selectAll('*').remove();     
        // Set the dimensions of the canvas / graph
        var margin = {top: 30, right: 120, bottom: 70, left: 50},
        width = this.divContainer.node().getBoundingClientRect().width - margin.left - margin.right,
        height = this.height - margin.top - margin.bottom;
        
        // set size of the legend
        var legendInfo = {left: width + margin.left, top: margin.top, spacing: 20}

        // Set the ranges
        var x = d3.scaleTime().range([0, width]);  
        var y = d3.scaleLinear().range([height, 0]);
        // Define the line
        var priceline = d3.line()	
                            .x(function(d) { return x(d.date); })
                            .y(function(d) { return y(d.value); });

        // Adds the svg canvas
        this.svg
            .style("width", "100%")
            .attr("height", this.height)
            .append("g")
            .attr("transform", 
            "translate(" + margin.left + "," + margin.top + ")");

        
        // Scale the range of the data
        x.domain(d3.extent(this.data, function(d) { return d.date; }));
        y.domain([0, d3.max(this.data, function(d) { return d.value; })]);
        
        // Nest the entries by symbol
        var dataNest = d3.nest()
            .key(function(d) {return d.name;})
            .entries(this.data);

        // Loop through each symbol / key
        dataNest.forEach(function(d,i) { 
            self.svg
                .append("path")
                .attr("data-legend", d.key)
                .attr("class", "line")
                .attr("d", priceline(d.values))
                .style("stroke", function() { // Add the colors dynamically
                    return d.color = d.values[0].background;
                })
                .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
            // Add the Legend
            var legend = self.svg.append("g")
                .attr("class", "legend")
                .attr("transform", "translate(" + legendInfo.left + "," + (legendInfo.top+ i * legendInfo.spacing) + ")");
            legend.append("rect")
                .attr("x", 0)
                .attr("y", 0)
                .attr("width", 10)
                .attr("height", 10)
                .style("fill", d.values[0].background);
            legend.append("text")
                .attr("x", 20) 
                .attr("y", 10)
                .attr("class", "legend")    // style the legend
                .text(d.key); 
        });

        // move the legend to the intended location

        // Add the legend
        // var legend = this.svg.append("g")
        //                 .attr("class","legend")
        //                 .attr("transform","translate(50,30)")
        //                 .style("font-size","12px")
        //                 .call(d3.legend)

        // setTimeout(function() { 
        //     legend
        //     .style("font-size","20px")
        //     .attr("data-style-padding",10)
        //     .call(d3.legend)
        // },1000)

        // Add the X Axis
        this.svg.append("g")
            .attr("class", "axis")
            .attr("transform", "translate(" + margin.left + "," + (height + margin.top) + ")")
            .call(d3.axisBottom(x));
        // Add the Y Axis
        this.svg.append("g")
            .attr("class", "axis")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
            .call(d3.axisLeft(y));
    }

    // TODO: add method to return data for a line that is the aggregate of all lines shown
    getTotalLine() {

    }
}
  