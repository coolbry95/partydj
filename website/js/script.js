var BaseURL = "https://linode.shellcode.in/"
var DatasetURL = document.baseURI.substring(0, document.baseURI.lastIndexOf("/") + 1) + "assets/";

function LoadList(dataset, idName, callback, attribute){
  _.orderBy(dataset, ["upvotes", "downvotes"], ["desc", "asc"]);
  var idName = idName === undefined ? "list" : idName;
  var template = "<li><img class='cover' src='COVER' alt='TRACK'/><div class='meta'><p class='track'>TRACK - ALBUM</p><p class='artist'>ARTIST</p></div><div class='control'><span class='count'>UPVOTES</span><img class='vote upvote' src='./assets/upvote.svg' alt='Upvote' data-track='TRACKID' onclick='VoteAction(this, true)' /><span class='count'>DOWNVOTES</span><img class='vote downvote' src='./assets/downvote.svg' alt='Downvote' data-track='TRACKID' onclick='VoteAction(this, false)' /></div></li>";
  
  document.getElementById(idName).innerHTML = "";
  
  if (dataset[0] !== undefined && dataset[0] !== null && dataset != ""){
    $("#track").html(dataset[0]["name"] + " - " + dataset[0]["albumname"]);
    $("#artist").html(dataset[0]["artists"][0]["name"]);
    $("#cover").attr({"alt": dataset[0]["name"], "src": dataset[0]["images"][2]["url"]});
    $("#count").html(dataset[0]["upvotes"] + " UP");
  }
  
  for (var i = 1; i < dataset.length; ++i){
    var itemTemp = template;
    var itemID = dataset[i]["ID"];
    var itemUpvote = dataset[i]["upvotes"];
    var itemDownvote = dataset[i]["downvotes"];
    var itemCover = dataset[i]["images"][2]["url"];
    var itemTrack = dataset[i]["name"];
    var itemAlbum = dataset[i]["albumname"];
    var itemArtists = dataset[i]["artists"];
    var itemDuration = dataset[i]["Duration"]
    
    itemTemp = itemTemp.replace(/TRACKID/g, itemID).replace(/UPVOTES/g, itemUpvote).replace(/DOWNVOTES/g, itemDownvote).replace(/COVER/g, itemCover).replace(/TRACK/g, itemTrack).replace(/ALBUM/g, itemAlbum).replace(/ARTIST/g, itemArtists[0]["name"]);
    
    document.getElementById(idName).innerHTML += itemTemp;
  }
}

function LoadJSON(queryURL, queryParam, queryMethod, storageKey, callback, attribute){
  var httpRequest = new XMLHttpRequest();
  var queryMethod = queryMethod === undefined ? "GET" : queryMethod;
  var storageKey = storageKey === undefined ? "temp" : storageKey;
  
  httpRequest.open(queryMethod, queryURL + (queryMethod === "GET" ? "?" + $.param(queryParam) : "") , true);
  //httpRequest.setRequestHeader("Access-Control-Allow-Origin", "*");
  httpRequest.onload = function(e){
    if (httpRequest.readyState === 4 && httpRequest.status === 200){
      var response = httpRequest.responseText === "" ? undefined : JSON.parse(httpRequest.responseText);
      if (response !== undefined && response !== null && response !== ""){
        localStorage.setItem(storageKey, httpRequest.responseText);
        localStorage.setItem("playlistid",response["playlistid"]);
      }
      
      if (callback !== undefined && callback !== null) callback(response, attribute);
    }
    else
      console.error(httpRequest.statusText);
  }
  // if (!queryURL.includes("file:///")) httpRequest.send(queryMethod === "GET" ? null : $.param(queryParam));
  if (!queryURL.includes("file:///"))
    if (queryMethod === "GET") httpRequest.send(null);
    else {
      var formData = new FormData();
      for (var key in queryParam) formData.append(key, queryParam[key]);
      httpRequest.send(formData);
    }
}

function VoteAction(el){
  console.log(el.getAttribute("track"));
}

function Initialize(dataset, attribute){
  var localUUID = localStorage.getItem("uuid");
  if (localUUID === null || localUUID === ""){
    identity().get(function(UUID){
      document.cookie = "uuid=" + UUID + "; max-age=86400;";
      localStorage.setItem("uuid", UUID);
    });
  }
  
  var localPoolID = localStorage.getItem("poolid");
  if (localPoolID !== null && localPoolID !== ""){
    $("#join").addClass("inactive");
    $("#pool").html("Pool: " + localPoolID);
    $("#quit, #container, #play").removeClass("inactive");
    
    var localDataset = localStorage.getItem("dataset");
    if ((dataset === undefined || dataset === null || dataset === "") 
      && localDataset !== null && localDataset !== ""){
      LoadList(JSON.parse(localDataset)["songheap"]);
    }
  }
  
  if (dataset !== undefined && dataset !== null && dataset !== ""){
    LoadList(dataset["songheap"]);
  }

}

function Finalize(dataset, attribute){
  Initialize();
  if (dataset === undefined || dataset === null || dataset === "" || (attribute["refresh"] !== undefined && attribute["refresh"] !== null && attribute["refresh"] === true)){
    var queryResult = LoadJSON(BaseURL + "getpool", {}, "GET", "dataset", Initialize, {});
    //var queryResult = LoadJSON(DatasetURL + "dataset.json", {}, "GET", "dataset", Initialize, {});
  }
}

function JoinPool(){
  var inputPoolID = document.getElementById("poolid").value;
  if (inputPoolID !== undefined && inputPoolID !== null && inputPoolID !== ""){
    localStorage.setItem("poolid", inputPoolID);
    var queryResult = LoadJSON(BaseURL + "join_pool", {"poolShortId": inputPoolID, "userId": localStorage.getItem("uuid")}, "POST", "temp", Finalize, {"refresh": true});
    //var queryResult = LoadJSON(DatasetURL + "dataset.json", {"poolShortId": inputPoolID, "userId": localStorage.getItem("uuid")}, "GET", "temp", Initialize, {});
  }
}

function QuitPool(){
  localStorage.removeItem("poolid");
  localStorage.removeItem("dataset");
  window.location.reload(true);
}

function VoteAction(elem, upvote){
  console.log(BaseURL + (upvote ? "upvote" : "downvote"));
  var queryResult = LoadJSON(BaseURL + (upvote ? "upvote/" : "downvote/") +  localStorage["playlistid"] + "/" + elem.getAttribute("data-track"), {"userId": localStorage.getItem("uuid")}, "POST", "temp", Finalize, {"refresh": true});
}

function SearchMusic(){
  var inputKeyword = document.getElementById("keyword").value;
  if (inputKeyword !== undefined && inputKeyword !== null && inputKeywordv !== ""){
    $("#result").removeClass("inactive");
    $("#list").addClass("inactive");
    var queryResult = LoadJSON(BaseURL + "search_for_songs", {"search_query": inputKeyword, "number_of_results": 10}, "POST", "temp", function(data, attr){ LoadList(data, "result"); }, {});
    //var queryResult = LoadJSON(DatasetURL + "dataset.json", {}, "GET", "temp", function(data, attr){ console.LoadList(data, "result"); }, {});
  }
}

Initialize();

$("#quit").click(QuitPool);
$("#quit").singletap(QuitPool);