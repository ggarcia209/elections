const { SearchRequest, SearchResult, SearchResponse } = require('./server_pb.js');
const { RankingsRequest, RankingsResult, RankingsResponse } = require('./server_pb.js');
const { LookupRequest, LookupResponse } = require('./server_pb.js');
const { YrTotalRequest, YrTotalResult, YrTotalResponse } = require('./server_pb.js');
const { GetIndvRequest, GetCmteRequest, GetCandRequest } = require('./server_pb.js');
const { Empty } = require('./server_pb.js');
const { ViewClient } = require('./server_grpc_web_pb.js');

var viewSvc = new ViewClient('http://localhost:8080');

// call on load functions for each page
let data = document.querySelector("#main").onload = load()
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
        case "/totals/":
            getYrTotals()
            break;
        case "/view-object/":
            getEntity()
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
        let entry = i + ".  " + rslt.getName() + " - " + rslt.getEmployer() + " - " +rslt.getCity() + ", " + rslt.getState()
        // let entry = i + ".  " + rslt.getName() + " - " + rslt.getCity() + ", " + rslt.getState()
        if (rslt.getBucket() !== "individuals") {
            entry = i + ".  " + rslt.getName() + " - " + rslt.getCity() + ", " + rslt.getState()
        }
        resultsString += "<li class='list-full-item'><p>"+ entry +"</p><span class='list-full-years'><ul class='years-list'>"
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
        resultsString +=   "<a class='rank-link' href='"+link+"'>" + i +".  " + r.getName() + " - " + r.getCity() + ", " + r.getState() + " - " + "$" + r.getAmount().toLocaleString() + "</a>";
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
        let entry = r.getName() + " - " + r.getCity() + ", " + r.getState()

        resultsString += "<li class='list-full-item'>";
        resultsString +=  "<p>" + i + ".  " + "<a class='list-full-link' href='"+link+"'>" + entry + ": </a> " + "$" + r.getAmount().toLocaleString() + "</p>";
        resultsString += "</li>"
        i++
    });
    document.querySelector("#rankings-list-full").innerHTML = resultsString
}
// END RANKINGS OPERATIONS

// YEARLY TOTAL OPERATIONS
function newYrTotalRequest(year, category, party) {
    let request = new YrTotalRequest();
    request.setUid("test007");
    request.setYear(year);
    request.setCategory(category);
    request.setParty(party);
    return request
}

function getYrTotals() {
    let params = (new URL(document.location)).searchParams;
    let year = params.get("year")
    let category = params.get("category")
    let req = newYrTotalRequest(year, category, "ALL")
    viewSvc.viewYrTotals(req, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }

       displayYrTotals(resp, year, category)
    })
    return 
}

function displayYrTotals(resp, year, category) {
    let i = 1
    let cats = {"rec": "Funds Received", "donor": "Funds Contributed", "exp": "Funds Expensed"}
    let ptys = {"ALL": "Overall", "DEM": "Democrat", "REP": "Republican", "IND": "Independent", "OTH": "Other", "UNK": "Unknown"}
    let yt = resp.getYearlytotalList()
    let resultsString = ""
    resultsString += year + " - " + cats[category];
    document.querySelector("#yt-list-hdr").innerHTML = resultsString
    resultsString = ""
    yt.forEach(function (r) {
        resultsString += "<li class='list-full-item'>";
        resultsString +=   i + ".  " + ptys[r.getParty()] + ": " + "$" + r.getTotal().toLocaleString();
        resultsString += "</li>";
        i++
    });

    // console.log("resultsString: ", resultsString)
    document.querySelector("#yt-list").innerHTML = resultsString
    return
}
// END YEARLY TOTAL OPERATIONS

// VIEW OBJECT OPERATIONS
function newGetIndvRequest(year, bucket, ID) {
    let request = new GetIndvRequest();
    let years = [year];
    request.setUid("test007");
    request.setObjectid(ID)
    request.setBucket(bucket);
    request.setYearsList(years);

    return request
}

function newGetCmteRequest(year, bucket, ID) {
    let request = new GetCmteRequest();
    let years = [year];
    request.setUid("test007");
    request.setObjectid(ID)
    request.setBucket(bucket);
    request.setYearsList(years);

    return request
}

function newGetCandRequest(year, bucket, ID) {
    let request = new GetCandRequest();
    let years = [year];
    request.setUid("test007");
    request.setObjectid(ID)
    request.setBucket(bucket);
    request.setYearsList(years);

    return request
}

function getEntity() {
    let params = (new URL(document.location)).searchParams;
    let year = params.get("year")
    let bucket = params.get("bucket")
    let objID = params.get("id")
    console.log(year, bucket, objID)
    

    switch(bucket){
        case "individuals":
            let ireq = newGetIndvRequest(year, bucket, objID)
            viewSvc.viewIndividual(ireq, {}, (err, resp) => {
                if (err !== null) {
                    console.log("error:")
                    console.log(err)
                    return
                }
                displayIndv(resp, year)
            })
            break;
        case "cmte_tx_data":
            let cmreq = newGetCmteRequest(year, bucket, objID)
            viewSvc.viewCommittee(cmreq, {}, (err, resp) => {
                if (err !== null) {
                    console.log("error:")
                    console.log(err)
                    return
                }
                displayCmte(resp, year)
            })
            break;
        case "committees":
            let cmmreq = newGetCmteRequest(year, "cmte_tx_data", objID)
            viewSvc.viewCommittee(cmmreq, {}, (err, resp) => {
                if (err !== null) {
                    console.log("error:")
                    console.log(err)
                    return
                }
                displayCmte(resp, year)
            })
            break;
        case "candidates":
            let cnreq = newGetCandRequest(year, bucket, objID)
            viewSvc.viewCandidate(cnreq, {}, (err, resp) => {
                if (err !== null) {
                    console.log("error:")
                    console.log(err)
                    return
                }
                displayCand(resp, year)
            })
            break;
    }
}

function displayIndv(resp, year) {
    console.log("individual...")
    let indv = resp.getIndividual()
    console.log("indv: ")
    let resultsString = ""
    resultsString += "<h1 class='header-about'>" + indv.getName() +  " - " + year + "</h1>";
    resultsString += "<h2 class='header-sub-about'> ID: " + indv.getId() + "</h3>";
    resultsString += "<h4 class='header-sub-about'>" + indv.getCity() + ", " + indv.getState() +  "</h4>"
    if (indv.getOccupation() !== "" && indv.getEmployer() !== "") {
        resultsString += "<h4 class='header-sub-about'>"+ indv.getOccupation() + ", " + indv.getEmployer() +  "</h4>"
    }

    resultsString += "<div class='list-full-div-obj'>"
    resultsString += "<ul class='list-full-obj'>"
    resultsString += "<li class='list-full-item-obj'>"
    resultsString += "<span class='list-full-years-obj'><ul class='years-list-obj'>"
    resp.getYearsList().forEach(function (y) {
        let link = "http://localhost:8081/view-object/?year="+y+"&bucket="+resp.getBucket()+"&id="+indv.getId()
        resultsString += "<li class='years-list-item-obj'><a class='list-full-link' href='"+link+"'>"+ y +"</a></li>";
    })
    resultsString += "</ul>"
    resultsString += "</span>"
    resultsString += "</li></ul></div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Total Incoming: $"+ indv.getTotalinamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Incoming Transactions: " + indv.getTotalintxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average: $" + indv.getAvgtxin().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Total Outgoing: $"+ indv.getTotaloutamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Outgoing Transactions: " + indv.getTotalouttxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average: $" + indv.getAvgtxout().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    document.querySelector("#obj-summary-div").innerHTML = resultsString
    document.querySelector("#obj-summary-div").style.display = "block"

    // get senders sorted list
    let senders = { amts: indv.getSendersamtList(), txs: indv.getSenderstxsMap() }
    let senIDs = new Array()
    let sAmts = {}
    senders.amts.forEach(function (r) {
        let id = r.getId() + ""
        sAmts[id] = r.getTotal()
        senIDs.push(id)
    });
    console.log("senIDs:")
    console.log(senIDs)

    let slReq = newLookupRequest(senIDs)
    viewSvc.lookupObjByID(slReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let sendersString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = sAmts[id]
            let txs = senders.txs.get(id)
            let avg = amt / txs
            sendersString += "<li class='list-full-item'>";
            sendersString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + r.getName() + "  - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")</p>";
            sendersString += "</li>"
            i++
        })
        document.querySelector("#senders-list").innerHTML = sendersString
        document.querySelector("#indv-senders").style.display = "block"
    })   

    // get recipients sorted list
    let recs = { amts: indv.getRecipientsamtList(), txs: indv.getRecipientstxsMap() }
    let recIDs = new Array()
    rAmts = {}
    recs.amts.forEach(function (r) {
        let id = r.getId() + ""
        rAmts[id] = r.getTotal()
        recIDs.push(id)
    });
    console.log("recIDs:")
    console.log(recIDs)

    let rlReq = newLookupRequest(recIDs)
    viewSvc.lookupObjByID(rlReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let recipientsString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = rAmts[id]
            let txs = recs.txs.get(id)
            let avg = amt / txs
            recipientsString += "<li class='list-full-item'>";
            recipientsString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")</p>";
            recipientsString += "</li>"
            i++
        })
        document.querySelector("#recipients-list").innerHTML = recipientsString
        document.querySelector("#indv-recipients").style.display = "block"
    })
}

function displayCmte(resp, year) {
    console.log("committee...")
    let resultsString = ""

    // Summary info
    let cmte = resp.getCommittee()
    resultsString += "<h1 class='header-about'><a class='title-link' href='https://www.fec.gov/data/committee/"+cmte.getId()+"/'>" + cmte.getName() + " - " + year + "</a></h1>";
    resultsString += "<h2 class='header-sub-about'>" + cmte.getParty() + "</h2>";
    resultsString += "<h2 class='header-sub-about'> ID: " + cmte.getId() + "</h2>";
    if (cmte.getCandid() !== "") {
        resultsString += "<h2 class='header-sub-about'> Candidate ID: " + cmte.getCandid() + "</h2>";
    }
    if (cmte.getConnectedorg() !== "") {
        resultsString += "<h2 class='header-sub-about'> Organization: " + cmte.getConnectedorg() + "</h2>";
    }
    resultsString += "<h4 class='header-sub-about'>" + cmte.getCity() + ", " + cmte.getState() +  "</h4>"

    resultsString += "<div class='list-full-div-obj'>"
    resultsString += "<ul class='list-full-obj'>"
    resultsString += "<li class='list-full-item-obj'>"
    resultsString += "<span class='list-full-years-obj'><ul class='years-list-obj'>"
    resp.getYearsList().forEach(function (y) {
        let link = "http://localhost:8081/view-object/?year="+y+"&bucket="+resp.getBucket()+"&id="+cmte.getId()
        resultsString += "<li class='years-list-item-obj'><a class='list-full-link' href='"+link+"'>"+ y +"</a></li>";
    })
    resultsString += "</ul>"
    resultsString += "</span>"
    resultsString += "</li></ul></div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    if (cmte.getDesignation() !== "") {
        resultsString += "<li class='list-view-item'><p>Designation: "+ cmte.getDesignation() + "</p></li>"
    }
    if (cmte.getType() !== "") {
        resultsString += "<li class='list-view-item'><p>Type: "+ cmte.getType() + "</p></li>"
    }
    if (cmte.getTresname() !== "") {
        resultsString += "<li class='list-view-item'><p>Treasurer: "+ cmte.getTresname() + "</p></li>"
    }
    resultsString += "</ul>"
    resultsString += "</div>"

    let txData = resp.getTxdata()

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Contributions Total: $"+ txData.getContributionsinamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Contributions Received: " + txData.getContributionsintxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average Contribution: $" + txData.getAvgcontributionin().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Other Receipts Total: $"+ txData.getOtherreceiptsinamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Other Receipts: " + txData.getOtherreceiptsintxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average Other Receipt: $" + txData.getAvgotherin().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Total Incoming: $"+ txData.getTotalincomingamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Total Incoming Transactions: " + txData.getTotalincomingtxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average Incoming Transaction: $" + txData.getAvgincoming().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Transfers to Other Committees Total: $"+ txData.getTransfersamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Transers to Other Committees: " + txData.getTransferstxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average Transfer: $" + txData.getAvgtransfer().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Expenditures Total: $"+ txData.getExpendituresamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Expenditure Transactions: " + txData.getExpenditurestxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average Expenditure: $" + txData.getAvgexpenditure().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Total Outgoing: $"+ txData.getTotaloutgoingamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Total Outgoing Transactions: " + txData.getTotaloutgoingtxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average Outgoing Transaction: $" + txData.getAvgoutgoing().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    document.querySelector("#obj-summary-div").innerHTML = resultsString
    document.querySelector("#obj-summary-div").style.display = "block"

    // Senders/Receivers
    // Top Individual Contributors
    let topIndv = { amts: txData.getTopindvcontributorsamtList(), txs: txData.getTopindvcontributorstxsMap() }
    let indvIDs = new Array()
    let iAmts = {}
    topIndv.amts.forEach(function (r) {
        let id = r.getId() + ""
        iAmts[id] = r.getTotal()
        indvIDs.push(id)
    });
    console.log("indvIDs:")
    console.log(indvIDs)

    let inReq = newLookupRequest(indvIDs)
    viewSvc.lookupObjByID(inReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let topIndvString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = iAmts[id]
            let txs = topIndv.txs.get(id)
            let avg = amt / txs
            let record = r.getName() + " - " + r.getEmployer() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            if (r.getEmployer() == "") {
                record = r.getName() + "  - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            }
            topIndvString += "<li class='list-full-item'>";
            topIndvString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + record + "</p>";
            topIndvString += "</li>"
            i++
        })
        document.querySelector("#top-indv-list").innerHTML = topIndvString
        document.querySelector("#cmte-top-indv").style.display = "block"
    })

    // Top Committee Contributors
    let topCmte = { amts: txData.getTopcmteorgcontributorsamtList(), txs: txData.getTopcmteorgcontributorstxsMap() }
    let cmteIDs = new Array()
    let cAmts = {}
    topCmte.amts.forEach(function (r) {
        let id = r.getId() + ""
        cAmts[id] = r.getTotal()
        cmteIDs.push(id)
    });
    console.log("cmteIDs:")
    console.log(cmteIDs)

    let cmReq = newLookupRequest(cmteIDs)
    viewSvc.lookupObjByID(cmReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let topCmteString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = cAmts[id]
            let txs = topCmte.txs.get(id)
            let avg = amt / txs
            topCmteString += "<li class='list-full-item'>";
            topCmteString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + r.getName() + "  - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")</p>";
            topCmteString += "</li>"
            i++
        })
        document.querySelector("#top-cmte-list").innerHTML = topCmteString
        document.querySelector("#cmte-top-cmte").style.display = "block"
    })

    // Transfer Recipients
    let tr = { amts: txData.getTransferrecsamtList(), txs: txData.getTransferrecstxsMap() }
    let trIDs = new Array()
    let trAmts = {}
    tr.amts.forEach(function (r) {
        let id = r.getId() + ""
        trAmts[id] = r.getTotal()
        trIDs.push(id)
    });
    console.log("recIDs:")
    console.log(trIDs)

    let trReq = newLookupRequest(trIDs)
    viewSvc.lookupObjByID(trReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let trRecsString  = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = trAmts[id]
            let txs = tr.txs.get(id)
            let avg = amt / txs
            trRecsString += "<li class='list-full-item'>";
            trRecsString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")</p>";
            trRecsString += "</li>"
            i++
        })
        document.querySelector("#tr-rec-list").innerHTML = trRecsString
        document.querySelector("#cmte-transfer-recs").style.display = "block"
    })

    // Top Expenditure Recipients
    let exps = { amts: txData.getTopexprecipientsamtList(), txs: txData.getTopexprecipientstxsMap() }
    let expIDs = new Array()
    let eAmts = {}
    exps.amts.forEach(function (r) {
        let id = r.getId() + ""
        eAmts[id] = r.getTotal()
        expIDs.push(id)
    });
    console.log("expIDs:")
    console.log(expIDs)

    let exReq = newLookupRequest(expIDs)
    viewSvc.lookupObjByID(exReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let expRecsString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = eAmts[id]
            let txs = exps.txs.get(id)
            let avg = amt / txs
            let entry = r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            if (bucket !== "individuals") {
                entry = entry = r.getName() + "  - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            }
            expRecsString += "<li class='list-full-item'>";
            expRecsString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + entry +"</p>";
            expRecsString += "</li>"
            i++
        })
        document.querySelector("#exp-rec-list").innerHTML = expRecsString
        document.querySelector("#cmte-exp-recs").style.display = "block"
    })
}

function displayCand(resp, year) {
    console.log("candidate...")
    let cand = resp.getCandidate()
    let resultsString = ""
    resultsString += "<h1 class='header-about'><a class='title-link' href='https://www.fec.gov/data/candidate/"+cand.getId()+"/'>" + cand.getName() + " - " + year + "</a></h1>";
    resultsString += "<h2 class='header-sub-about'>" + cand.getParty() + " - " + cand.getOffice() + "</h2>"
    resultsString += "<h2 class='header-sub-about'> ID: " + cand.getId() + "</h2>"
    resultsString += "<h2 class='header-sub-about'> PCC: " + cand.getPcc() + "</h2>"
    resultsString += "<h4 class='header-sub-about'>" + cand.getCity() + ", " + cand.getState() +  "</h4>"

    resultsString += "<div class='list-full-div-obj'>"
    resultsString += "<ul class='list-full-obj'>"
    resultsString += "<li class='list-full-item-obj'>"
    resultsString += "<span class='list-full-years-obj'><ul class='years-list-obj'>"
    resp.getYearsList().forEach(function (y) {
        let link = "http://localhost:8081/view-object/?year="+y+"&bucket="+resp.getBucket()+"&id="+cand.getId()
        resultsString += "<li class='years-list-item-obj'><a class='list-full-link' href='"+link+"'>"+ y +"</a></li>";
    })
    resultsString += "</ul>"
    resultsString += "</span>"
    resultsString += "</li></ul></div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Total Incoming: $"+ cand.getTotaldirectinamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Incoming Transactions: " + cand.getTotaldirectintxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average: $" + cand.getAvgdirectin().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    resultsString += "<div class='list-full-div'>"
    resultsString += "<ul class='list-full'>"
    resultsString += "<li class='list-view-item'><p>Total Outgoing: $"+ cand.getTotaldirectoutamt().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Outgoing Transactions: " + cand.getTotaldirectouttxs().toLocaleString() + "</p></li>"
    resultsString += "<li class='list-view-item'><p>Average: $" + cand.getAvgdirectout().toLocaleString() + "</p></li>"
    resultsString += "</ul>"
    resultsString += "</div>"

    document.querySelector("#obj-summary-div").innerHTML = resultsString
    document.querySelector("#obj-summary-div").style.display = "block"

    // get senders sorted list
    let senders = { amts: cand.getDirectsendersamtsList(), txs: cand.getDirectsenderstxsMap() }
    let senIDs = new Array()
    let sAmts = {}
    senders.amts.forEach(function (r) {
        let id = r.getId() + ""
        sAmts[id] = r.getTotal()
        senIDs.push(id)
    });
    console.log("senIDs:")
    console.log(senIDs)

    let slReq = newLookupRequest(senIDs)
    viewSvc.lookupObjByID(slReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        let res = resp.getResultsList()
        let sendersString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = sAmts[id]
            let txs = senders.txs.get(id)
            let avg = amt / txs
            let entry = r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            if (bucket !== "individuals") {
                entry = entry = r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            }
            sendersString += "<li class='list-full-item'>";
            sendersString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + entry +"</p>";
            sendersString += "</li>"
            i++
        })
        document.querySelector("#senders-list").innerHTML = sendersString
        document.querySelector("#indv-senders").style.display = "block"
    })   

    // get recipients sorted list
    let recs = { amts: cand.getDirectrecipientsamtsList(), txs: cand.getDirectrecipientstxsMap() }
    let recIDs = new Array()
    let rAmts = {}
    recs.amts.forEach(function (r) {
        let id = r.getId() + ""
        rAmts[id] = r.getTotal()
        recIDs.push(id)
    });
    console.log("recIDs:")
    console.log(recIDs)

    let rlReq = newLookupRequest(recIDs)
    viewSvc.lookupObjByID(rlReq, {}, (err, resp) => {
        if (err !== null) {
            console.log("error:")
            console.log(err)
            return
        }
        console.log(resp)
        let res = resp.getResultsList()
        let recipientsString = ""
        let i = 1
        res.forEach(function (r) {
            let id = r.getId() + ""
            let bucket = r.getBucket()
            let link = "http://localhost:8081/view-object/?year="+year+"&bucket="+bucket+"&id="+id
            let amt = rAmts[id]
            let txs = recs.txs.get(id)
            let avg = amt / txs
            let entry = r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            if (bucket !== "individuals") {
                entry = entry = r.getName() + " - " + r.getCity() + ", " + r.getState() +"</a>: $" + amt.toLocaleString() + " (Avg: $"+avg.toLocaleString()+")"
            }
            recipientsString += "<li class='list-full-item'>";
            recipientsString +=   "<p>"+i +". <a class='list-full-link' href='"+link+"'>" + entry + "</p>";
            recipientsString += "</li>"
            i++
        })
        document.querySelector("#recipients-list").innerHTML = recipientsString
        document.querySelector("#indv-recipients").style.display = "block"
    })
}


function newLookupRequest(reqIDs) {
    let req = new LookupRequest()
    req.setUid("test007")
    req.setObjectidsList(reqIDs)
    return req
}

// END VIEW OBJECT OPERATIONS

// HELPER FUNCTIONS
function createNode(ele) {
    return document.createElement(ele);
}

function append(par, chi) {
    par.append(chi)
}