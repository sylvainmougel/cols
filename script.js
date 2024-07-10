const list = document.getElementById("stash");
let button = document.getElementById('button');

function showProfile() {
    let altitude = []
    let distance = []
    let traces = [];
    let tag = document.getElementById('tag').value;
    const client = algoliasearch('1QMZVCS1V5', 'cb6b989a18ef9a3070d5b5a54001a3da');
    const index = client.initIndex('profiles');
    const bar = new Promise((resolve, reject) => {
        index.search("", {filters: tag}).then(({hits}) => {
            console.log(hits)
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

// Add a 'click' event listener to the button
button.addEventListener('click', function () {
    showProfile()
});



