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

//	NodeID         NodeID
//	LeaderID       uint64
//	ViewID         uint64
//	SystemNodes    []uint64
//	LastUpdateTime string


let nodeID;
let leaderID;
let viewID;
let systemNodes;
let lastUpdateTime;
let records;
let totalFiles;

function updateStatus() {
  fetch('/api/status')
  .then(response => response.json())
  .then(json => {
   nodeID = json.NodeID;
   leaderID = json.LeaderID;
   viewID = json.ViewID;
   systemNodes = json.SystemNodes;
   lastUpdateTime = json.LastUpdateTime;
   records = json.Records;
   totalFiles = json.TotalFiles;

   document.getElementById("NodeID").innerHTML = nodeID;
   document.getElementById("LeaderID").innerHTML = leaderID;
 
   document.getElementById("ViewID").innerHTML = viewID;   
   document.getElementById("SystemNodes").innerHTML = systemNodes;   
   document.getElementById("LastUpdateTime").innerHTML = lastUpdateTime;

   document.getElementById("Records").innerHTML = records;
   document.getElementById("TotalFiles").innerHTML = totalFiles;

   document.getElementById("connectionIndicator").classList.remove("hidden");
   document.getElementById("errorIndicator").classList.add("hidden");

   })
 .catch((error) => {

  document.getElementById("connectionIndicator").classList.add("hidden");
  document.getElementById("errorIndicator").classList.remove("hidden");

   console.error('Error:', error);
 });

}

let updateStatusInterval = setInterval(updateStatus, 1000);

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
select.replaceChildren();

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
