const { SearchRequest, SearchResult, SearchResponse } = require('./server_pb.js');
const { RankingsRequest, RankingsResult, RankingsResponse } = require('./server_pb.js');
const { LookupRequest, LookupResponse } = require('./server_pb.js');
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
        case "/rankings/":
            loadRankingsNew()
            break;
        case "/rankings-list/":
            getRankings()
            break;
    }
}

// SEARCH OPERATIONS
function search() {
    let params = (new URL(document.location)).searchParams;
    let query = params.get("q");
    getSearchResponse(query)
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
    let i = 1
    let resultsString = "";
        res.forEach(function (r) {
        rslt = new SearchResult()
        rslt = r
        let entry = i + ".  " + rslt.getName() + " - " +rslt.getCity() + ", " + rslt.getState()
        resultsString += "<li class='list-full-item'><p>"+entry+"</p><span class='list-full-years'><ul class='years-list'>"
        rslt.getYearsList().forEach(function (y) {
            let link = "http://localhost:8081/view-object/?year="+y+"&bucket="+rslt.getBucket()+"&id="+rslt.getId()
            resultsString += "<li class='years-list-item'><a class='list-full-link' href='"+link+"'>"+ y +"</a></li>";
        })
        resultsString += "</ul>"
        resultsString += "</span>";
        resultsString += "</li>";
        i++
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

function loadRankingsNew() {
    // clearRankingsLists()
    let params = (new URL(document.location)).searchParams;
    let year = params.get("year");
    let bktCat = params.get("category").split("-");
    let bucket = bktCat[0]
    let category = bktCat[1]

    getRankingsPrev("#ranks-page-1", year, bucket, category, "ALL")
    getRankingsPrev("#ranks-page-2", year, bucket, category, "DEM")
    getRankingsPrev("#ranks-page-3", year, bucket, category, "REP")
    getRankingsPrev("#ranks-page-4", year, bucket, category, "IND")
    getRankingsPrev("#ranks-page-5", year, bucket, category, "OTH")
    getRankingsPrev("#ranks-page-6", year, bucket, category, "UNK")
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

function getRankingsPrev(ulId, year, bucket, category, party) {
    party = party + "-pre"
    let req = newRankingsRequest(year, bucket, category, party)
    viewSvc.viewRankings(req, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        if (ulId !== "#ranks-pre-main-1" && ulId !== "#ranks-pre-main-2") { // exclude home page
            console.log("if !ulId ", ulId)
            displayRankingsTitle(ulId, year, bucket, category, party)
            displayRankingsBtn(ulId, year, bucket, category, party)
        }  
        displayRankings(resp, ulId, year, bucket)
    })
    return 
}

function displayRankingsTitle(ulID, year, bucket, category, party) {
    let titleIDDict = {
        "#ranks-page-1": "#ranks-title-1",
        "#ranks-page-2": "#ranks-title-2",
        "#ranks-page-3": "#ranks-title-3",
        "#ranks-page-4": "#ranks-title-4",
        "#ranks-page-5": "#ranks-title-5",
        "#ranks-page-6": "#ranks-title-6",
    }
    let titleDict = {
        "individuals-donor": "Individual Contributors",
        "individuals-rec": "Individual Recipients",
        "cmte_tx_data-rec": "Committtee Recipients",
        "cmte_tx_data-donor": "Committee Donors",
        "cmte_tx_data-exp": "Committee Spenders",
        "candidates-rec": "Candidate Recipients",
        "candidates-donor": "Candidate Donors",
        "candidates-exp": "Candidate Spenders",
    }
    let partyDict = {
        "ALL-pre": "Overall",
        "DEM-pre": "Democrat",
        "REP-pre": "Republican",
        "IND-pre": "Independent",
        "OTH-pre": "Other",
        "UNK-pre": "Unknown",
    }
    let titleID = titleIDDict[ulID]
    let hdr = document.querySelector(titleID) 
    let pty = partyDict[party]
    console.log("party: " + pty)
    let title = titleDict[bucket + "-" + category]
    let titleString = year + " - " + title + " - " + pty
    hdr.innerHTML = titleString
}

function displayRankings(resp, ulId, year, bucket) {
    console.log(ulId)
    let i = 1
    let rnk = new RankingsResult()
    rnk = resp.getRankings()
    let rnkList = rnk.getRankingslistList()
    let resultsString = ""
    rnkList.forEach(function (r) {
        let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+r.getId()
        resultsString += "<li class='rank-item'>";
        resultsString +=   "<a class='rank-link' href='"+link+"'>" + i +")  " + r.getName() + " - " + r.getCity() + ", " + r.getState() + " - " + "$" + r.getAmount() + "</a>";
        resultsString += "</li>"
        i++
    });
    if (i == 1) {
        hideEmptyRankings(ulId)
        return
    }
    // console.log("resultsString: ", resultsString)
    document.querySelector(ulId).innerHTML = resultsString
    return
}

function displayRankingsBtn(ulId, year, bucket, category, party) {
    let btnIds = {
        "#ranks-page-1": "#ranks-btn-1",
        "#ranks-page-2": "#ranks-btn-2",
        "#ranks-page-3": "#ranks-btn-3",
        "#ranks-page-4": "#ranks-btn-4",
        "#ranks-page-5": "#ranks-btn-5",
        "#ranks-page-6": "#ranks-btn-6",
    }
    let partyDict = {
        "ALL-pre": "ALL",
        "DEM-pre": "DEM",
        "REP-pre": "REP",
        "IND-pre": "IND",
        "OTH-pre": "OTH",
        "UNK-pre": "UNK",
    }
    let id = btnIds[ulId]
    console.log("ID: "+id)
    let link = document.querySelector(id)
    console.log("link: "+link)
    let pty = partyDict[party]
    link.href = "http://localhost:8081/rankings-list/?year="+year+"&bucket="+bucket+"&category="+category+"&party="+pty
}

function hideEmptyRankings(ulId) {
    let IdDict = {
        "#ranks-page-1": "#ranks-single-1",
        "#ranks-page-2": "#ranks-single-2",
        "#ranks-page-3": "#ranks-single-3",
        "#ranks-page-4": "#ranks-single-4",
        "#ranks-page-5": "#ranks-single-5",
        "#ranks-page-6": "#ranks-single-6",
    }
    let ID = IdDict[ulId]
    let div = document.querySelector(ID) 
    div.style.display = "none"
}

function getRankings() {
    let params = (new URL(document.location)).searchParams;
    let year = params.get("year")
    let bucket = params.get("bucket")
    let category = params.get("category")
    let party = params.get("party")
    console.log(year, bucket, category, party)
    let req = newRankingsRequest(year, bucket, category, party)
    viewSvc.viewRankings(req, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        displayRankingsFullTitle(year, bucket, category, party)
        displayRankingsAll(resp, year, bucket)
    })
    return 
}

function displayRankingsFullTitle(year, bucket, category, party) {
    let titleDict = {
        "individuals-donor": "Individual Contributors",
        "individuals-rec": "Individual Recipients",
        "cmte_tx_data-rec": "Committtee Recipients",
        "cmte_tx_data-donor": "Committee Donors",
        "cmte_tx_data-exp": "Committee Spenders",
        "candidates-rec": "Candidate Recipients",
        "candidates-donor": "Candidate Donors",
        "candidates-exp": "Candidate Spenders",
    }
    let partyDict = {
        "ALL": "Overall",
        "DEM": "Democrat",
        "REP": "Republican",
        "IND": "Independent",
        "OTH": "Other",
        "UNK": "Unknown",
    }
    let hdr = document.querySelector("#rankings-list-title")
    let title = titleDict[bucket + "-" + category]
    let pty = partyDict[party]
    let titleString = year + " - " + title + " - " + pty
    hdr.innerHTML = titleString
}

function displayRankingsAll(resp, year, bucket) {
    let i = 1
    let rnk = new RankingsResult()
    rnk = resp.getRankings()
    let rnkList = rnk.getRankingslistList()
    let resultsString = ""
    rnkList.forEach(function (r) {
        let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+r.getId()
        resultsString += "<li class='list-full-item'>";
        resultsString +=   "<a class='list-full-link' href='"+link+"'>" + r.getName() + " - " + r.getCity() + ", " + r.getState() + " - " + "$" + r.getAmount() + "</a>";
        resultsString += "</li>"
        i++
    });
    document.querySelector("#rankings-list-full").innerHTML = resultsString
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