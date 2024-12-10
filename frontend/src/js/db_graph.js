var sqlTable_list = {}


window.onload = function() {
        print_graph();
        // Sample data
       /* let data = [
            {category: 'A', value: 30},
            {category: 'B', value: 80},
            {category: 'C', value: 45},
            {category: 'D', value: 60},
            {category: 'E', value: 20},
            {category: 'F', value: 90},
            {category: 'G', value: 55}
        ];
        barplot_v2(data)*/
};


function print_graph() {

    let table_name = "chronologie"
    let rows_limit = "1000"
    parent.window.go.main.App.Get_header_table(table_name, rows_limit).then(resultat => {
        console.log("graph:")
        console.log(resultat)

/*
        let data = [
            {category: 'A', value: 50},
            {category: 'B', value: 50},
            {category: 'C', value: 50},
            {category: 'D', value: 50},
            {category: 'E', value: 50},
            {category: 'F', value: 50},
            {category: 'G', value: 50}
        ];
        barplot_1(data)*/
        d3.csv("https://raw.githubusercontent.com/holtzy/data_to_viz/master/Example_dataset/1_OneNum.csv").then( function(data) {
            barplot_2(data);
        });


    }).catch(err => console.error("Error db_info:", err));

}



/*
let data = [
    {category: 'A', value: 50},
    {category: 'B', value: 50},
    {category: 'C', value: 50},
    {category: 'D', value: 50},
    {category: 'E', value: 50},
    {category: 'F', value: 50},
    {category: 'G', value: 50}
];
barplot_1(data)
*/
function barplot_1(data,
                   width=400,
                   height=150,
                   margin = {top: 20, right: 30, bottom: 20, left: 40},
                   bg_color = "lightblue",
                   bar_color = "#69b3a2",
                   x_max = 1000,
                   ticks_number = 100,
                   title="Title",
                   x_label="X Axis",
                   y_label="Y Axis") {
    // HTML container
    let balise_name = "graph"
    let grapth_balise = document.getElementById(balise_name);
    grapth_balise.innerHTML = ""
    grapth_balise.style = "background-color: "+bg_color+";" +
        "width:"+width+"px;" +
        "height:"+height+"px;" +
        "padding: 15px "+ (margin.right+30) +"px " + (margin.bottom + 10) +"px 15px"





    // SVG container
    const svg = d3.select("#"+balise_name).append("svg")
        .attr("width", width)
        .attr("height", height+margin.bottom);

    // Scales
    const x = d3.scaleBand()
        .domain(data.map(d => d.category))
        .range([margin.left, width - margin.right])
        .padding(0.1);

    const y = d3.scaleLinear()
        .domain([0, d3.max(data, d => d.value)]).nice()
        .range([height - margin.bottom, margin.top]);

    // Bars
    svg.selectAll(".bar")
        .data(data)
        .enter().append("rect")
        .attr("class", "bar")
        .attr("x", d => x(d.category))
        .attr("y", d => y(d.value))
        .attr("width", x.bandwidth())
        .attr("height", d => y(0) - y(d.value))
        .attr("fill", "steelblue");

    // X-axis
    svg.append("g")
        .attr("transform", `translate(0,${height - margin.bottom})`)
        .call(d3.axisBottom(x));

    // Y-axis
    svg.append("g")
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y));

    // Axis labels
    svg.append("text")
        .attr("x", width / 2)
        .attr("y", (height + margin.bottom/2) )
        .attr("text-anchor", "middle")
        .text("Category");


    svg.append("text")
        .attr("transform", "rotate(-90)")
        .attr("x", -height / 2)
        .attr("y", 15)
        .attr("text-anchor", "middle")
        .text("Value");



}



/*
d3.csv("https://raw.githubusercontent.com/holtzy/data_to_viz/master/Example_dataset/1_OneNum.csv").then( function(data) {
    barplot_2(data);
});
*/
function barplot_2(data,
                    balise="my_dataviz",
                    column="price",
                    //aggregate="sum", // [sum, mean, count]
                    width=400,
                    height=150,
                    margin = {top: 20, right: 30, bottom: 40, left: 40},
                    bg_color = "lightblue",
                    bar_color = "#69b3a2",
                    x_max = 1000,
                    ticks_number = 100,
                    title="Title",
                    x_label="X Axis",
                    y_label="Y Axis"){
    if (x_label === "X Axis") {x_label=column}

    // HTML container
    let grapth_balise = document.getElementById(balise);
    grapth_balise.innerHTML = ""
    grapth_balise.style = "background-color: "+bg_color+";" +
        "width:"+width+"px;" +
        "height:"+height+"px;" +
        "padding: 15px "+ (margin.right+30) +"px " + (margin.bottom + 20) +"px 15px"

    // append the svg object to the body of the page
    const svg = d3.select("#"+balise)
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            `translate(${margin.left}, ${margin.top})`);

    // Title
    svg.append("text")
        .attr("x", (width / 2))
        .attr("y", 5 - (margin.top / 2))
        .attr("text-anchor", "middle")
        .style("font-size", "16px")
        .style("text-decoration", "underline")
        .text(title);

    // X axis: scale and draw:
    const x = d3.scaleLinear()
        .domain([0, x_max])
        .range([0, width]);
    svg.append("g")
        .attr("transform", `translate(0, ${height})`)
        .call(d3.axisBottom(x));

    // X Axis label
    svg.append("text")
        .attr("transform",
            `translate(${width / 2} ,${height + margin.top + 15})`)
        .style("text-anchor", "middle")
        .text(x_label);


    // set the parameters for the histogram
    const histogram = d3.histogram()
        .value(function(d) { return d[column]; })   // I need to give the vector of value
        .domain(x.domain())
        .thresholds(x.ticks(ticks_number));

    // And apply this function to data to get the bins
    const bins = histogram(data);

    // Y axis: scale and draw:
    const y = d3.scaleLinear()
        .range([height, 0]);
    y.domain([0, d3.max(bins, function(d) { return d.length; })]);   // d3.hist has to be called before the Y axis obviously
    svg.append("g")
        .call(d3.axisLeft(y));

    // Y Axis label
    svg.append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", -5 - margin.left)
        .attr("x", 0 - (height / 2))
        .attr("dy", "1em")
        .style("text-anchor", "middle")
        .text(y_label);


    // Add a tooltip div. Here I define the general feature of the tooltip: stuff that do not depend on the data point.
    // Its opacity is set to 0: we don't see it by default.
    const tooltip = d3.select("#"+balise)
        .append("div")
        .style("opacity", 0)
        .attr("class", "tooltip")
        .style("background-color", "black")
        .style("color", "white")
        .style("border-radius", "5px")
        .style("padding", "10px")
        .style("position", "absolute");

    // A function that change this tooltip when the user hover a point.
    // Its opacity is set to 1: we can now see it. Plus it set the text and position of tooltip depending on the datapoint (d)
    const showTooltip = function(event,d) {
        console.log(d)
        tooltip
            .transition()
            .duration(100)
            .style("opacity", 1)
        tooltip
            .html("Range: " + d.x0 + " - " + d.x1) // tooltip text
            .style("left", (event.pageX + 10) + "px") // Adjust tooltip position
            .style("top", (event.pageY - 30) + "px"); // Adjust tooltip position
        }
    const moveTooltip = function(event,d) {
        tooltip
            .style("left", (event.pageX + 10) + "px") // Adjust tooltip position
            .style("top", (event.pageY - 30) + "px"); // Adjust tooltip position
        }
    // A function that change this tooltip when the leaves a point: just need to set opacity to 0 again
    const hideTooltip = function(event,d) {
        tooltip
            .transition()
            .duration(100)
            .style("opacity", 0)
    }

    // append the bar rectangles to the svg element
    svg.selectAll("rect")
        .data(bins)
        .join("rect")
        .attr("x", 1)
        .attr("transform", function(d) { return `translate(${x(d.x0)}, ${y(d.length)})`})
        .attr("width", function(d) { return Math.max(0, x(d.x1) - x(d.x0) - 1); })
        .attr("height", function(d) { return height - y(d.length); })
        .style("fill", bar_color)
        // Show tooltip on hover
        .on("mouseover", showTooltip )
        .on("mousemove", moveTooltip )
        .on("mouseleave", hideTooltip )
}


//data ={"values": [32.2, 33.4, 34.5, 35.8, 37.0, 38.2, 39.3, 40.3, 41.2, 42.0, 42.7, 43.2, 43.7, 44.1, 44.6, 45.2, 45.9, 46.6, 47.3, 48.1, 49.0, 49.9, 50.8, 51.7, 52.6, 53.7, 54.7, 55.9, 57.0, 58.2, 59.3, 60.2, 61.1, 61.9, 62.6, 63.6, 64.6, 65.2, 65.0, 64.0, 61.8, 58.6, 54.2, 48.1, 41.2, 34.6, 29.3, 26.1, 24.9, 25.6, 28.3, 33.3, 39.6, 46.0, 51.3, 54.4, 55.9, 55.7, 53.9, 50.2, 45.5, 40.6, 36.4, 33.9, 32.5, 32.2, 33.0, 35.0, 37.7, 40.7, 43.4, 45.2, 46.6, 47.5, 48.0, 47.8, 47.4, 47.0, 46.9, 47.6, 48.7, 50.5, 52.7, 56.1, 59.7, 63.1, 65.7, 66.9, 66.9, 65.5, 62.9, 58.8, 53.8, 48.7, 44.2, 40.9, 38.5, 37.0, 36.3, 36.3, 36.8, 37.7, 38.7, 39.4, 40.2, 41.2, 42.2, 43.7, 45.2, 46.6, 47.5, 47.7, 47.3, 46.3, 44.1, 41.4, 38.6, 36.2, 34.5, 34.0]}
function dateplot_3(data,
                   balise="my_graphDate",
                   column="price",
                   //aggregate="sum", // [sum, mean, count]
                   width=400,
                   height=150,
                   margin = {top: 20, right: 30, bottom: 40, left: 40},
                   bg_color = "lightblue",
                   bar_color = "#69b3a2",
                   x_max = 1000,
                   ticks_number = 100,
                   title="Title",
                   x_label="X Axis",
                   y_label="Y Axis"){
    if (x_label === "X Axis") {x_label=column}

    // HTML container
    let grapth_balise = document.getElementById(balise);
    grapth_balise.innerHTML = ""
    grapth_balise.style = "background-color: "+bg_color+";" +
        "width:"+width+"px;" +
        "height:"+height+"px;" +
        "padding: 15px "+ (margin.right+30) +"px " + (margin.bottom + 20) +"px 15px"

    // append the svg object to the body of the page
    const svg = d3.select("#"+balise)
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            `translate(${margin.left}, ${margin.top})`);



}

