const tagIDs = ['tdf2024', 's11', 's14','s15', 's19']
let elems = new Map()
// Create the document
for (let tag of tagIDs) {
    var checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.id = tag;
    checkbox.name = tag;
    checkbox.addEventListener('click', function () {
        showProfile()
    });

    var label = document.createElement('label');
    label.htmlFor = tag;
    label.appendChild(document.createTextNode("#" + tag));

    var container = document.getElementById('filters');
    container.appendChild(checkbox);
    container.appendChild(label);
    elems.set(tag, checkbox)
    if (tag.includes("tdf")) {
        var br = document.createElement('br');
        container.appendChild(br);
    }
}


function showProfile() {
    // Compute filter
    let filters = ""
    for (const [key, value] of elems) {
        if (value.checked) {
            if (filters !== "") {
                filters += " OR "
            }
            filters += key
        }
    }
    console.log("filters: ", filters)
    if (filters=== "") {
        let plot = document.getElementById('plot');
        plot.innerHTML = '';
        return
    }
    // Display graph.
    let altitude = []
    let distance = []
    let traces = [];
    const client = algoliasearch('1QMZVCS1V5', 'cb6b989a18ef9a3070d5b5a54001a3da');
    const index = client.initIndex('profiles');
    const bar = new Promise((resolve, reject) => {
        index.search("", {filters: filters}).then(({hits}) => {
            hits.forEach((res => {
                altitude = res["slaltitude"]
                distance = res["distance"]
                slope = res["slope"]
                colors = []
                slope.forEach(s => {
                    if (s <= 6) {
                        colors.push("green")
                    } else if (s <= 8) {
                        colors.push("blue")
                    } else if (s <= 10) {
                        colors.push("red")
                    } else {
                        colors.push("black")
                    }
                });
                const trace = {
                    x: distance,
                    y: altitude,
                    mode: 'markers+lines',
                    marker: {'color': colors},
                    line: {'color': 'gray'},
                    type: 'scatter',
                    name: res["name"]
                };
                traces.push(trace)
            }))
            resolve()
        });
    })
    bar.then(() => {
            Plotly.newPlot('plot', traces);
        }
    )

}

showProfile()







