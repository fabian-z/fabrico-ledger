// Import our custom CSS
import '../scss/styles.scss'

// Import all of Bootstrap's JS
import * as bootstrap from 'bootstrap'

import '@fortawesome/fontawesome-free/js/fontawesome'
import '@fortawesome/fontawesome-free/js/solid'
import '@fortawesome/fontawesome-free/js/regular'
import '@fortawesome/fontawesome-free/js/brands'

import Chart from 'chart.js/auto'

Chart.defaults.backgroundColor = 'transparent';
Chart.defaults.borderColor = '#6b6b6b';
Chart.defaults.color = '#fff';

// Fetch node status / id

let nodeID;
fetch('/api/status')
 .then(response => response.json())
 .then(json => {
	nodeID = json.NodeID;
	document.getElementById("NodeID").innerHTML = nodeID;   
	})
.catch((error) => {
  console.error('Error:', error);
});

document.getElementById("upload-submit").addEventListener('click', (event) => {
// TODO error handling

	fetch("api/addfile", {
    body: new FormData(document.getElementById("upload-form")),
    method: "post",
});

});

document.getElementById("fabricate-submit").addEventListener('click', (event) => {
// TODO error handling
	fetch("api/fabricate", {
    body: new FormData(document.getElementById("fabricate-form")),
    method: "post",
});

});

document.getElementById("fabricate-refresh").addEventListener('click', refreshFabricate);
document.getElementById("nav-fabricate").addEventListener('click', refreshFabricate);

function refreshFabricate() {
	
fetch('/api/availabledata')
 .then(response => response.json())
 .then(json => {


let select = document.getElementById('partSelect');


for (const data of json){
   // data.OriginatingID
   // data.FileHash
   // data.Remaining
    let opt = document.createElement('option');
    opt.value = data.FileHash;
    opt.innerHTML = "Part " + data.FileHash.substr(0, 8) + " from Node " + data.OriginatingID + " (" + data.Remaining + " available)";
    select.appendChild(opt);
}

	})
.catch((error) => {
  console.error('Error:', error);
});
}

// (Dummy) System information

// Graphs
const ctx = document.getElementById('myChart')
// eslint-disable-next-line no-unused-vars
const myChart = new Chart(ctx, {
  type: 'line',
  data: {
    labels: [
      'Sunday',
      'Monday',
      'Tuesday',
      'Wednesday',
      'Thursday',
      'Friday',
      'Saturday'
    ],
    datasets: [{
		label: 'Fabrications',
      data: [
        15,
        21,
        18,
        24,
        23,
        240,
        120
      ],
      lineTension: 0,
    },
    {
		label: 'Uploads',
      data: [
        5,
        15,
        83,
        5,
        12,
        24,
        42
      ],
      lineTension: 0,
    }]
  },
  options: {
    scales: {
      yAxes: [{
        ticks: {
          beginAtZero: false
        }
      }]
    },
    legend: {
      display: false
    }
  }
})
