// code from https://gist.github.com/d3noob/f8b7f107ba25c21971851728520224cb

import * as d3 from 'd3';
import d3Tip from '@lix/d3-tip';
d3.tip = d3Tip;

export class MultilineGraph {
    constructor(height, data, id) {
        this.height = height;
        this.data = this.parseData(data);

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

        this.redraw();

        window.addEventListener('resize', function(context) {
            return function() {
                context.redraw()
            };
        }(this));
    }

    redraw() {
        self = this;
        this.svg.selectAll('*').remove();        
        // Set the dimensions of the canvas / graph
        var margin = {top: 30, right: 100, bottom: 70, left: 50},
        width = this.divContainer.node().getBoundingClientRect().width - margin.left - margin.right,
        height = this.height - margin.top - margin.bottom;
        
        // set size of the legend
        var legendInfo = {left: width + margin.left, top: margin.top, spacing: 20}

        // Set the ranges
        var x = d3.scaleTime().range([0, width]);  
        var y = d3.scaleLinear().range([height, 0]);
        // Define the line
        var priceline = d3.line()	
                            .x(function(d) { return x(d.Date); })
                            .y(function(d) { return y(d.Value); });

        // Adds the svg canvas
        this.svg
            .style("width", "100%")
            .attr("height", this.height)
            .append("g")
            .attr("transform", 
            "translate(" + margin.left + "," + margin.top + ")");

        
        // Scale the range of the data
        x.domain(d3.extent(this.data, function(d) { return d.Date; }));
        y.domain([0, d3.max(this.data, function(d) { return d.Value; })]);
        // Nest the entries by symbol
        var dataNest = d3.nest()
            .key(function(d) {return d.Name;})
            .entries(this.data);
        // set the colour scale
        var color = d3.scaleOrdinal(d3.schemeCategory10);

        var legendSpace = width/dataNest.length; // spacing for the legend
        // Loop through each symbol / key
        dataNest.forEach(function(d,i) { 
            self.svg.append("path")
                .attr("class", "line")
                .style("stroke", function() { // Add the colours dynamically
                    return d.color = d.values[0].Background; })
                .attr("d", priceline(d.values))
                .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
            // Add the Legend
            self.svg.append("text")
                .attr("x", legendInfo.left)  // space legend
                .attr("y", legendInfo.top + i * legendInfo.spacing)
                .attr("class", "legend")    // style the legend
                .style("fill", function() { // Add the colours dynamically
                    return d.color = d.values[0].Background; })
                .text(d.key); 
        });
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

    parseData(dataStr) {
        var data = JSON.parse(dataStr);

        data.forEach(function(d) {
            d.Date = Date.parse(d.Date);
        });
        return data;
    }

    // TODO: add method to return data for a line that is the aggregate of all lines shown
    getTotalLine() {

    }
}
  