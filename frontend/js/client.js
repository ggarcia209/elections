const { SearchRequest, SearchResult, SearchResponse } = require('./server_pb.js');
const { RankingsRequest, RankingsResult, RankingsResponse } = require('./server_pb.js');
const { YrTotalRequest, YrTotalResult, YrTotalResponse } = require('./server_pb.js');
const { GetObjRequest, GetObjResponse } = require('./server_pb.js');
const { Empty } = require('./server_pb.js');
const { ViewClient } = require('./server_grpc_web_pb.js');
grpc.web = require('grpc-web');

var viewSvc = new ViewClient('http://localhost:8080');
const resultsList = document.getElementById("search-list")


function getQueryVariable(variable) {
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split("=");
        if (pair[0] === variable) {
            return decodeURIComponent(pair[1].replace(/\+/g, "%20"));
        }
    }
}

function newSearchRequest(text)  {
    let request = new SearchRequest;
    request.setUID = "test007";
    request.setText = text;
    return request
}
  
function getSearchResponse(query) {
    // let text = document.getElementById("search-main").value;
    let req = newSearchRequest(query);
    let resp = new SearchResponse
    resp = viewSvc.SearchQuery(req);
    return resp
}

function displaySearchResults(resp) {
    let res = resp.getResults()
    let li = createNode('li'),
        a = createNode('a');
    li.setAttribute('class', 'list-full-item')
    a.setAttribute('href', '#')
    a.innerHTML = "test test test"
    append(li, a)
    append(resultsList, li)
    // resultPages from previous example
    let resultsString = "";
        res.forEach(function (r) {
        rslt = new SearchResult
        rslt = r
        resultsString += "<li class='list-full-item'>";
        resultsString +=   "<a class='list-full-link' href='#'>" + result.getName() + "</a>";
        resultsString += "</li>"
    });
    document.querySelector("#search-results-main").innerHTML = resultsString;
}

function newRankingsRequest(year, bucket, category, party) {
    let request = new RankingsRequest;
    request.setUID = "test007";
    request.setYear = year;
    request.setBucket = bucket;
    request.setCategory = category;
    request.setParty = party;
    return request
}

function getRankings(year, bucket, category, party) {
    let req = newRankingsRequest(year, bucket, category, party)
    let resp = viewSvc.viewRankings(req)
    return resp
}

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