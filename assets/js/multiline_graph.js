import * as d3 from 'd3';
import d3Tip from '@lix/d3-tip';
d3.tip = d3Tip;

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

        this.createGraph();
        // self = this;
        // $(window).resize(function() {
        //     self.chart.width(self.getWidth(self.margin))
        //         .height(self.getHeight(self.margin));
    
        //     self.svg
        //         .attr('width', self.getWidth(self.margin))
        //         .attr('height', self.getHeight(self.margin))
        //         .call(self.chart);
        // });
    }

    generateData() {
        var d = `symbol,date,price
        MSFT,Jan 2000,39.81
        MSFT,Feb 2000,36.35
        MSFT,Mar 2000,43.22
        MSFT,Apr 2000,28.37
        MSFT,May 2000,25.45
        MSFT,Jun 2000,32.54
        MSFT,Jul 2000,28.4
        MSFT,Aug 2000,28.4
        MSFT,Sep 2000,24.53
        MSFT,Oct 2000,28.02
        MSFT,Nov 2000,23.34
        MSFT,Dec 2000,17.65
        MSFT,Jan 2001,24.84
        MSFT,Feb 2001,24
        MSFT,Mar 2001,22.25
        MSFT,Apr 2001,27.56
        MSFT,May 2001,28.14
        MSFT,Jun 2001,29.7
        MSFT,Jul 2001,26.93
        MSFT,Aug 2001,23.21
        MSFT,Sep 2001,20.82
        MSFT,Oct 2001,23.65
        MSFT,Nov 2001,26.12
        MSFT,Dec 2001,26.95
        MSFT,Jan 2002,25.92
        MSFT,Feb 2002,23.73
        MSFT,Mar 2002,24.53
        MSFT,Apr 2002,21.26
        MSFT,May 2002,20.71
        MSFT,Jun 2002,22.25
        MSFT,Jul 2002,19.52
        MSFT,Aug 2002,19.97
        MSFT,Sep 2002,17.79
        MSFT,Oct 2002,21.75
        MSFT,Nov 2002,23.46
        MSFT,Dec 2002,21.03
        MSFT,Jan 2003,19.31
        MSFT,Feb 2003,19.34
        MSFT,Mar 2003,19.76
        MSFT,Apr 2003,20.87
        MSFT,May 2003,20.09
        MSFT,Jun 2003,20.93
        MSFT,Jul 2003,21.56
        MSFT,Aug 2003,21.65
        MSFT,Sep 2003,22.69
        MSFT,Oct 2003,21.45
        MSFT,Nov 2003,21.1
        MSFT,Dec 2003,22.46
        MSFT,Jan 2004,22.69
        MSFT,Feb 2004,21.77
        MSFT,Mar 2004,20.46
        MSFT,Apr 2004,21.45
        MSFT,May 2004,21.53
        MSFT,Jun 2004,23.44
        MSFT,Jul 2004,23.38
        MSFT,Aug 2004,22.47
        MSFT,Sep 2004,22.76
        MSFT,Oct 2004,23.02
        MSFT,Nov 2004,24.6
        MSFT,Dec 2004,24.52
        MSFT,Jan 2005,24.11
        MSFT,Feb 2005,23.15
        MSFT,Mar 2005,22.24
        MSFT,Apr 2005,23.28
        MSFT,May 2005,23.82
        MSFT,Jun 2005,22.93
        MSFT,Jul 2005,23.64
        MSFT,Aug 2005,25.35
        MSFT,Sep 2005,23.83
        MSFT,Oct 2005,23.8
        MSFT,Nov 2005,25.71
        MSFT,Dec 2005,24.29
        MSFT,Jan 2006,26.14
        MSFT,Feb 2006,25.04
        MSFT,Mar 2006,25.36
        MSFT,Apr 2006,22.5
        MSFT,May 2006,21.19
        MSFT,Jun 2006,21.8
        MSFT,Jul 2006,22.51
        MSFT,Aug 2006,24.13
        MSFT,Sep 2006,25.68
        MSFT,Oct 2006,26.96
        MSFT,Nov 2006,27.66
        MSFT,Dec 2006,28.13
        MSFT,Jan 2007,29.07
        MSFT,Feb 2007,26.63
        MSFT,Mar 2007,26.35
        MSFT,Apr 2007,28.3
        MSFT,May 2007,29.11
        MSFT,Jun 2007,27.95
        MSFT,Jul 2007,27.5
        MSFT,Aug 2007,27.34
        MSFT,Sep 2007,28.04
        MSFT,Oct 2007,35.03
        MSFT,Nov 2007,32.09
        MSFT,Dec 2007,34
        MSFT,Jan 2008,31.13
        MSFT,Feb 2008,26.07
        MSFT,Mar 2008,27.21
        MSFT,Apr 2008,27.34
        MSFT,May 2008,27.25
        MSFT,Jun 2008,26.47
        MSFT,Jul 2008,24.75
        MSFT,Aug 2008,26.36
        MSFT,Sep 2008,25.78
        MSFT,Oct 2008,21.57
        MSFT,Nov 2008,19.66
        MSFT,Dec 2008,18.91
        MSFT,Jan 2009,16.63
        MSFT,Feb 2009,15.81
        MSFT,Mar 2009,17.99
        MSFT,Apr 2009,19.84
        MSFT,May 2009,20.59
        MSFT,Jun 2009,23.42
        MSFT,Jul 2009,23.18
        MSFT,Aug 2009,24.43
        MSFT,Sep 2009,25.49
        MSFT,Oct 2009,27.48
        MSFT,Nov 2009,29.27
        MSFT,Dec 2009,30.34
        MSFT,Jan 2010,28.05
        MSFT,Feb 2010,28.67
        MSFT,Mar 2010,28.8
        AMZN,Jan 2000,64.56
        AMZN,Feb 2000,68.87
        AMZN,Mar 2000,67
        AMZN,Apr 2000,55.19
        AMZN,May 2000,48.31
        AMZN,Jun 2000,36.31
        AMZN,Jul 2000,30.12
        AMZN,Aug 2000,41.5
        AMZN,Sep 2000,38.44
        AMZN,Oct 2000,36.62
        AMZN,Nov 2000,24.69
        AMZN,Dec 2000,15.56
        AMZN,Jan 2001,17.31
        AMZN,Feb 2001,10.19
        AMZN,Mar 2001,10.23
        AMZN,Apr 2001,15.78
        AMZN,May 2001,16.69
        AMZN,Jun 2001,14.15
        AMZN,Jul 2001,12.49
        AMZN,Aug 2001,8.94
        AMZN,Sep 2001,5.97
        AMZN,Oct 2001,6.98
        AMZN,Nov 2001,11.32
        AMZN,Dec 2001,10.82
        AMZN,Jan 2002,14.19
        AMZN,Feb 2002,14.1
        AMZN,Mar 2002,14.3
        AMZN,Apr 2002,16.69
        AMZN,May 2002,18.23
        AMZN,Jun 2002,16.25
        AMZN,Jul 2002,14.45
        AMZN,Aug 2002,14.94
        AMZN,Sep 2002,15.93
        AMZN,Oct 2002,19.36
        AMZN,Nov 2002,23.35
        AMZN,Dec 2002,18.89
        AMZN,Jan 2003,21.85
        AMZN,Feb 2003,22.01
        AMZN,Mar 2003,26.03
        AMZN,Apr 2003,28.69
        AMZN,May 2003,35.89
        AMZN,Jun 2003,36.32
        AMZN,Jul 2003,41.64
        AMZN,Aug 2003,46.32
        AMZN,Sep 2003,48.43
        AMZN,Oct 2003,54.43
        AMZN,Nov 2003,53.97
        AMZN,Dec 2003,52.62
        AMZN,Jan 2004,50.4
        AMZN,Feb 2004,43.01
        AMZN,Mar 2004,43.28
        AMZN,Apr 2004,43.6
        AMZN,May 2004,48.5
        AMZN,Jun 2004,54.4
        AMZN,Jul 2004,38.92
        AMZN,Aug 2004,38.14
        AMZN,Sep 2004,40.86
        AMZN,Oct 2004,34.13
        AMZN,Nov 2004,39.68
        AMZN,Dec 2004,44.29
        AMZN,Jan 2005,43.22
        AMZN,Feb 2005,35.18
        AMZN,Mar 2005,34.27
        AMZN,Apr 2005,32.36
        AMZN,May 2005,35.51
        AMZN,Jun 2005,33.09
        AMZN,Jul 2005,45.15
        AMZN,Aug 2005,42.7
        AMZN,Sep 2005,45.3
        AMZN,Oct 2005,39.86
        AMZN,Nov 2005,48.46
        AMZN,Dec 2005,47.15
        AMZN,Jan 2006,44.82
        AMZN,Feb 2006,37.44
        AMZN,Mar 2006,36.53
        AMZN,Apr 2006,35.21
        AMZN,May 2006,34.61
        AMZN,Jun 2006,38.68
        AMZN,Jul 2006,26.89
        AMZN,Aug 2006,30.83
        AMZN,Sep 2006,32.12
        AMZN,Oct 2006,38.09
        AMZN,Nov 2006,40.34
        AMZN,Dec 2006,39.46
        AMZN,Jan 2007,37.67
        AMZN,Feb 2007,39.14
        AMZN,Mar 2007,39.79
        AMZN,Apr 2007,61.33
        AMZN,May 2007,69.14
        AMZN,Jun 2007,68.41
        AMZN,Jul 2007,78.54
        AMZN,Aug 2007,79.91
        AMZN,Sep 2007,93.15
        AMZN,Oct 2007,89.15
        AMZN,Nov 2007,90.56
        AMZN,Dec 2007,92.64
        AMZN,Jan 2008,77.7
        AMZN,Feb 2008,64.47
        AMZN,Mar 2008,71.3
        AMZN,Apr 2008,78.63
        AMZN,May 2008,81.62
        AMZN,Jun 2008,73.33
        AMZN,Jul 2008,76.34
        AMZN,Aug 2008,80.81
        AMZN,Sep 2008,72.76
        AMZN,Oct 2008,57.24
        AMZN,Nov 2008,42.7
        AMZN,Dec 2008,51.28
        AMZN,Jan 2009,58.82
        AMZN,Feb 2009,64.79
        AMZN,Mar 2009,73.44
        AMZN,Apr 2009,80.52
        AMZN,May 2009,77.99
        AMZN,Jun 2009,83.66
        AMZN,Jul 2009,85.76
        AMZN,Aug 2009,81.19
        AMZN,Sep 2009,93.36
        AMZN,Oct 2009,118.81
        AMZN,Nov 2009,135.91
        AMZN,Dec 2009,134.52
        AMZN,Jan 2010,125.41
        AMZN,Feb 2010,118.4
        AMZN,Mar 2010,128.82
        IBM,Jan 2000,100.52
        IBM,Feb 2000,92.11
        IBM,Mar 2000,106.11
        IBM,Apr 2000,99.95
        IBM,May 2000,96.31
        IBM,Jun 2000,98.33
        IBM,Jul 2000,100.74
        IBM,Aug 2000,118.62
        IBM,Sep 2000,101.19
        IBM,Oct 2000,88.5
        IBM,Nov 2000,84.12
        IBM,Dec 2000,76.47
        IBM,Jan 2001,100.76
        IBM,Feb 2001,89.98
        IBM,Mar 2001,86.63
        IBM,Apr 2001,103.7
        IBM,May 2001,100.82
        IBM,Jun 2001,102.35
        IBM,Jul 2001,94.87
        IBM,Aug 2001,90.25
        IBM,Sep 2001,82.82
        IBM,Oct 2001,97.58
        IBM,Nov 2001,104.5
        IBM,Dec 2001,109.36
        IBM,Jan 2002,97.54
        IBM,Feb 2002,88.82
        IBM,Mar 2002,94.15
        IBM,Apr 2002,75.82
        IBM,May 2002,72.97
        IBM,Jun 2002,65.31
        IBM,Jul 2002,63.86
        IBM,Aug 2002,68.52
        IBM,Sep 2002,53.01
        IBM,Oct 2002,71.76
        IBM,Nov 2002,79.16
        IBM,Dec 2002,70.58
        IBM,Jan 2003,71.22
        IBM,Feb 2003,71.13
        IBM,Mar 2003,71.57
        IBM,Apr 2003,77.47
        IBM,May 2003,80.48
        IBM,Jun 2003,75.42
        IBM,Jul 2003,74.28
        IBM,Aug 2003,75.12
        IBM,Sep 2003,80.91
        IBM,Oct 2003,81.96
        IBM,Nov 2003,83.08
        IBM,Dec 2003,85.05
        IBM,Jan 2004,91.06
        IBM,Feb 2004,88.7
        IBM,Mar 2004,84.41
        IBM,Apr 2004,81.04
        IBM,May 2004,81.59
        IBM,Jun 2004,81.19
        IBM,Jul 2004,80.19
        IBM,Aug 2004,78.17
        IBM,Sep 2004,79.13
        IBM,Oct 2004,82.84
        IBM,Nov 2004,87.15
        IBM,Dec 2004,91.16
        IBM,Jan 2005,86.39
        IBM,Feb 2005,85.78
        IBM,Mar 2005,84.66
        IBM,Apr 2005,70.77
        IBM,May 2005,70.18
        IBM,Jun 2005,68.93
        IBM,Jul 2005,77.53
        IBM,Aug 2005,75.07
        IBM,Sep 2005,74.7
        IBM,Oct 2005,76.25
        IBM,Nov 2005,82.98
        IBM,Dec 2005,76.73
        IBM,Jan 2006,75.89
        IBM,Feb 2006,75.09
        IBM,Mar 2006,77.17
        IBM,Apr 2006,77.05
        IBM,May 2006,75.04
        IBM,Jun 2006,72.15
        IBM,Jul 2006,72.7
        IBM,Aug 2006,76.35
        IBM,Sep 2006,77.26
        IBM,Oct 2006,87.06
        IBM,Nov 2006,86.95
        IBM,Dec 2006,91.9
        IBM,Jan 2007,93.79
        IBM,Feb 2007,88.18
        IBM,Mar 2007,89.44
        IBM,Apr 2007,96.98
        IBM,May 2007,101.54
        IBM,Jun 2007,100.25
        IBM,Jul 2007,105.4
        IBM,Aug 2007,111.54
        IBM,Sep 2007,112.6
        IBM,Oct 2007,111
        IBM,Nov 2007,100.9
        IBM,Dec 2007,103.7
        IBM,Jan 2008,102.75
        IBM,Feb 2008,109.64
        IBM,Mar 2008,110.87
        IBM,Apr 2008,116.23
        IBM,May 2008,125.14
        IBM,Jun 2008,114.6
        IBM,Jul 2008,123.74
        IBM,Aug 2008,118.16
        IBM,Sep 2008,113.53
        IBM,Oct 2008,90.24
        IBM,Nov 2008,79.65
        IBM,Dec 2008,82.15
        IBM,Jan 2009,89.46
        IBM,Feb 2009,90.32
        IBM,Mar 2009,95.09
        IBM,Apr 2009,101.29
        IBM,May 2009,104.85
        IBM,Jun 2009,103.01
        IBM,Jul 2009,116.34
        IBM,Aug 2009,117
        IBM,Sep 2009,118.55
        IBM,Oct 2009,119.54
        IBM,Nov 2009,125.79
        IBM,Dec 2009,130.32
        IBM,Jan 2010,121.85
        IBM,Feb 2010,127.16
        IBM,Mar 2010,125.55
        AAPL,Jan 2000,25.94
        AAPL,Feb 2000,28.66
        AAPL,Mar 2000,33.95
        AAPL,Apr 2000,31.01
        AAPL,May 2000,21
        AAPL,Jun 2000,26.19
        AAPL,Jul 2000,25.41
        AAPL,Aug 2000,30.47
        AAPL,Sep 2000,12.88
        AAPL,Oct 2000,9.78
        AAPL,Nov 2000,8.25
        AAPL,Dec 2000,7.44
        AAPL,Jan 2001,10.81
        AAPL,Feb 2001,9.12
        AAPL,Mar 2001,11.03
        AAPL,Apr 2001,12.74
        AAPL,May 2001,9.98
        AAPL,Jun 2001,11.62
        AAPL,Jul 2001,9.4
        AAPL,Aug 2001,9.27
        AAPL,Sep 2001,7.76
        AAPL,Oct 2001,8.78
        AAPL,Nov 2001,10.65
        AAPL,Dec 2001,10.95
        AAPL,Jan 2002,12.36
        AAPL,Feb 2002,10.85
        AAPL,Mar 2002,11.84
        AAPL,Apr 2002,12.14
        AAPL,May 2002,11.65
        AAPL,Jun 2002,8.86
        AAPL,Jul 2002,7.63
        AAPL,Aug 2002,7.38
        AAPL,Sep 2002,7.25
        AAPL,Oct 2002,8.03
        AAPL,Nov 2002,7.75
        AAPL,Dec 2002,7.16
        AAPL,Jan 2003,7.18
        AAPL,Feb 2003,7.51
        AAPL,Mar 2003,7.07
        AAPL,Apr 2003,7.11
        AAPL,May 2003,8.98
        AAPL,Jun 2003,9.53
        AAPL,Jul 2003,10.54
        AAPL,Aug 2003,11.31
        AAPL,Sep 2003,10.36
        AAPL,Oct 2003,11.44
        AAPL,Nov 2003,10.45
        AAPL,Dec 2003,10.69
        AAPL,Jan 2004,11.28
        AAPL,Feb 2004,11.96
        AAPL,Mar 2004,13.52
        AAPL,Apr 2004,12.89
        AAPL,May 2004,14.03
        AAPL,Jun 2004,16.27
        AAPL,Jul 2004,16.17
        AAPL,Aug 2004,17.25
        AAPL,Sep 2004,19.38
        AAPL,Oct 2004,26.2
        AAPL,Nov 2004,33.53
        AAPL,Dec 2004,32.2
        AAPL,Jan 2005,38.45
        AAPL,Feb 2005,44.86
        AAPL,Mar 2005,41.67
        AAPL,Apr 2005,36.06
        AAPL,May 2005,39.76
        AAPL,Jun 2005,36.81
        AAPL,Jul 2005,42.65
        AAPL,Aug 2005,46.89
        AAPL,Sep 2005,53.61
        AAPL,Oct 2005,57.59
        AAPL,Nov 2005,67.82
        AAPL,Dec 2005,71.89
        AAPL,Jan 2006,75.51
        AAPL,Feb 2006,68.49
        AAPL,Mar 2006,62.72
        AAPL,Apr 2006,70.39
        AAPL,May 2006,59.77
        AAPL,Jun 2006,57.27
        AAPL,Jul 2006,67.96
        AAPL,Aug 2006,67.85
        AAPL,Sep 2006,76.98
        AAPL,Oct 2006,81.08
        AAPL,Nov 2006,91.66
        AAPL,Dec 2006,84.84
        AAPL,Jan 2007,85.73
        AAPL,Feb 2007,84.61
        AAPL,Mar 2007,92.91
        AAPL,Apr 2007,99.8
        AAPL,May 2007,121.19
        AAPL,Jun 2007,122.04
        AAPL,Jul 2007,131.76
        AAPL,Aug 2007,138.48
        AAPL,Sep 2007,153.47
        AAPL,Oct 2007,189.95
        AAPL,Nov 2007,182.22
        AAPL,Dec 2007,198.08
        AAPL,Jan 2008,135.36
        AAPL,Feb 2008,125.02
        AAPL,Mar 2008,143.5
        AAPL,Apr 2008,173.95
        AAPL,May 2008,188.75
        AAPL,Jun 2008,167.44
        AAPL,Jul 2008,158.95
        AAPL,Aug 2008,169.53
        AAPL,Sep 2008,113.66
        AAPL,Oct 2008,107.59
        AAPL,Nov 2008,92.67
        AAPL,Dec 2008,85.35
        AAPL,Jan 2009,90.13
        AAPL,Feb 2009,89.31
        AAPL,Mar 2009,105.12
        AAPL,Apr 2009,125.83
        AAPL,May 2009,135.81
        AAPL,Jun 2009,142.43
        AAPL,Jul 2009,163.39
        AAPL,Aug 2009,168.21
        AAPL,Sep 2009,185.35
        AAPL,Oct 2009,188.5
        AAPL,Nov 2009,199.91
        AAPL,Dec 2009,210.73
        AAPL,Jan 2010,192.06
        AAPL,Feb 2010,204.62
        AAPL,Mar 2010,223.02`
        
        return d3.csvParse(d);
    }

    createGraph() {
        self = this;
        // Set the dimensions of the canvas / graph
        var margin = {top: 30, right: 20, bottom: 70, left: 50},
        width = this.divContainer.node().getBoundingClientRect().width - margin.left - margin.right,
        height = this.height - margin.top - margin.bottom;
        // Parse the date / time
        var parseDate = d3.timeParse("%b %Y");
        // Set the ranges
        var x = d3.scaleTime().range([0, width]);  
        var y = d3.scaleLinear().range([height, 0]);
        // Define the line
        var priceline = d3.line()	
                            .x(function(d) { return x(d.date); })
                            .y(function(d) { return y(d.price); });

        // Adds the svg canvas
        this.svg.attr("width", width + margin.left + margin.right)
            .attr("height", this.height)
            .append("g")
            .attr("transform", 
                "translate(" + margin.left + "," + margin.top + ")");

        this.data.forEach(function(d) {
            d.date = parseDate(d.date);
            d.price = +d.price;
        });
        // Scale the range of the data
        x.domain(d3.extent(this.data, function(d) { return d.date; }));
        y.domain([0, d3.max(this.data, function(d) { return d.price; })]);
        // Nest the entries by symbol
        var dataNest = d3.nest()
            .key(function(d) {return d.symbol;})
            .entries(this.data);
        // set the colour scale
        var color = d3.scaleOrdinal(d3.schemeCategory10);
        var legendSpace = width/dataNest.length; // spacing for the legend
        // Loop through each symbol / key
        dataNest.forEach(function(d,i) { 
            self.svg.append("path")
                .attr("class", "line")
                .style("stroke", function() { // Add the colours dynamically
                    return d.color = color(d.key); })
                .attr("d", priceline(d.values))
                .attr("transform", "translate(" + margin.left + ",0)");
            // Add the Legend
            self.svg.append("text")
                .attr("x", (legendSpace/2)+i*legendSpace)  // space legend
                .attr("y", height + (margin.bottom/2)+ 5)
                .attr("class", "legend")    // style the legend
                .attr("transform", "translate(" + margin.left + ",0)")
                .style("fill", function() { // Add the colours dynamically
                    return d.color = color(d.key); })
                .text(d.key); 
        });
        // Add the X Axis
        this.svg.append("g")
            .attr("class", "axis")
            .attr("transform", "translate(" + margin.left + "," + height + ")")
            .call(d3.axisBottom(x));
        // Add the Y Axis
        this.svg.append("g")
            .attr("class", "axis")
            .attr("transform", "translate(" + margin.left + ",0)")
            .call(d3.axisLeft(y));
    }
}
  