const { SearchRequest, SearchResult, SearchResponse } = require('./server_pb.js');
const { RankingsRequest, RankingsResult, RankingsResponse } = require('./server_pb.js');
const { YrTotalRequest, YrTotalResult, YrTotalResponse } = require('./server_pb.js');
const { GetObjRequest, GetObjResponse } = require('./server_pb.js');
const { Empty } = require('./server_pb.js');
const { ViewClient } = require('./server_grpc_web_pb.js');

var viewSvc = new ViewClient('http://localhost:8080');

// call on load functions for each page
let data = document.getElementById("#main").onload = load()
function load() {
    let path = window.location.pathname
    console.log(path)
    switch(path){
        case "/":
            loadRankingsMain()
            break;
        case "/search-results/":
            search()
            break;
    }
}

// SEARCH OPERATIONS
function search() {
    let params = (new URL(document.location)).searchParams;
    let query = params.get("q");
    resp = getSearchResponse(query)
    console.log("displaying results...")
    displaySearchResults(resp)
}

function newSearchRequest(text)  {
    let request = new SearchRequest();
    request.setUid("test007")
    request.setText(text)
    return request
}
  
function getSearchResponse(query) {
    let request = newSearchRequest(query);
    viewSvc.searchQuery(request, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err.message)
            if (err.message == "MAX_LENGTH") {
                displayErrorMsg()
            }
            return
        }
        let msg = resp.getMsg()
        console.log(msg)
        if (msg == "NO_RESULTS") {
            displayNoResults()
            return
        }
        displaySearchResults(resp)
    });
    return
}

function displaySearchResults(resp) {
    let res = resp.getResultsList()
    let resultsString = "";
        res.forEach(function (r) {
        rslt = new SearchResult()
        rslt = r
        resultsString += "<li class='list-full-item'>";
        resultsString +=   "<a class='list-full-link' href='#'>" + rslt.getName() + " - " +rslt.getCity() + ", " + rslt.getState() + "</a>";
        resultsString += "</li>"
    });
    document.querySelector("#search-list").innerHTML = resultsString;
}

function displayNoResults() {
    let resultsString = "";
    resultsString += "<li class='list-full-item'>"
    resultsString += "No results found! Please use the search bar above to try again.";
    resultsString += "</li>";
    document.querySelector("#search-list").innerHTML = resultsString;
}

function displayErrorMsg() {
    let resultsString = "";
    resultsString += "<li class='list-full-item'>"
    resultsString += "Too many results! Please refine your search by adding one or more search terms and use the search bar above to try again.";
    resultsString += "</li>";
    document.querySelector("#search-list").innerHTML = resultsString;
}
// END SEARCH OPERATIONS

// RANKINGS OPRATIONS
function loadRankingsMain() {
    getRankingsPrev("#ranks-pre-main-1", "2020", "individuals", "donor", "ALL")
    getRankingsPrev("#ranks-pre-main-2", "2020", "cmte_tx_data", "rec", "ALL")
}


function newRankingsRequest(year, bucket, category, party) {
    let request = new RankingsRequest();
    request.setUid("test007")
    request.setYear(year)
    request.setBucket(bucket)
    request.setCategory(category)
    request.setParty(party)
    return request
}

function getRankings(ulId, year, bucket, category, party) {
    let req = newRankingsRequest(year, bucket, category, party)
    viewSvc.viewRankings(req, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        displayRankings(resp, ulId)
    })
    return 
}

function getRankingsPrev(ulId, year, bucket, category, party) {
    party = party + "-pre"
    let req = newRankingsRequest(year, bucket, category, party)
    viewSvc.viewRankings(req, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        console.log(resp)
        displayRankings(resp, ulId)
    })
    return 
}

function displayRankings(resp, ulId) {
    console.log(ulId)
    let i = 1
    let rnk = new RankingsResult()
    rnk = resp.getRankings()
    let rnkList = rnk.getRankingslistMap()
    let resultsString = "";
        rnkList.forEach(function (value, key) {
        resultsString += "<li class='rank-item'>";
        resultsString +=   "<a class='rank-link' href='#'>" + i +")  " + key + " - " + value + "</a>";
        resultsString += "</li>"
        i++
    });
    document.querySelector(ulId).innerHTML = resultsString;
}
// END RANKINGS OPERATIONS

// YEARLY TOTAL OPERATIONS
function newYrTotalRequest(year, category, party) {
    let request = new YrTotalRequest;
    request.setUID = "test007";
    request.setYear = year;
    request.setCategory = category;
    request.setParty = party;
    return request
}

function getYrTotal(year, category, party) {
    let req = newYrTotalRequest(year, category, party);
    let resp = viewSvc.viewYrTotals(req);
    return resp
}
// END YEARLY TOTAL OPERATIONS


// HELPER FUNCTIONS
function createNode(ele) {
    return document.createElement(ele);
}

function append(par, chi) {
    par.append(chi)
}