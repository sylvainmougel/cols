const list = document.getElementById("stash");
let button = document.getElementById('button');

function showProfile() {
    let altitude = []
    let distance = []
    let traces = [];
    const client = algoliasearch('1QMZVCS1V5', 'cb6b989a18ef9a3070d5b5a54001a3da');
    const index = client.initIndex('profiles');
    const bar = new Promise((resolve, reject) => {
        index.getObjects(['17275870', '5211636']).then(({results}) => {
            results.forEach((res => {
                altitude = res["profile"]["altitude"]["data"]
                distance = res["profile"]["distance"]["data"]
                const trace = {
                    x: distance,
                    y: altitude,
                    type: 'scatter'
                };
                traces.push(trace)
            }))
            resolve()
        });
        index.getObjects(['17275870', '5211636']).then(({res} )=> {
            console.log(res)
            res.forEach((hit => {

                const trace = {
                    x: distance,
                    y: altitude,
                    type: 'scatter'
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



