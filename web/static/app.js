// web/static/app.js
console.log("Chart is", typeof Chart, Chart);
function renderCharts(data) {
    const labels = [];
    const scores = [];
    const contrib = [];
    const bg = [];

    Object.entries(data).forEach(([nodeId, node]) => {
        const score = typeof node.score === 'number' ? node.score.toFixed(2) : "0.00";
        const contribution = typeof node.contribution === 'number' ? node.contribution.toFixed(2) : "0.00";

        labels.push(nodeId);
        scores.push(score);
        contrib.push(contribution);
        bg.push(node.is_anchor ? 'rgba(255,99,132,0.7)' : 'rgba(54,162,235,0.7)');
    });

    if (window.scoreChart && typeof window.scoreChart.destroy === 'function' && window.scoreChart instanceof Chart) {
        window.scoreChart.destroy();
    }


    window.scoreChart = new Chart(document.getElementById("scoreChart"), {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: "节点得分",
                data: scores,
                backgroundColor: bg
            }]
        },
        options: {
            plugins: {
                legend: { display: false },
                tooltip: {
                    callbacks: {
                        label: (ctx) => `得分: ${ctx.raw}`
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });

    window.contribChart = new Chart(document.getElementById("contribChart"), {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: "贡献值",
                data: contrib,
                backgroundColor: 'rgba(75, 192, 192, 0.6)'
            }]
        },
        options: {
            plugins: {
                legend: { display: false }
            },
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });
}

function fetchAndRender() {
    fetch("/list")
        .then(res => res.json())
        .then(data => renderCharts(data))
        .catch(err => console.error("Fetch error:", err));
}

fetchAndRender();
setInterval(fetchAndRender, 5000);