var BaseURL = "https://linode.shellcode.in/"
var DatasetURL = document.baseURI.substring(0, document.baseURI.lastIndexOf("/") + 1) + "assets/";

function LoadJSON(queryURL, queryParam, queryMethod, storageKey, callback, attribute){
  var httpRequest = new XMLHttpRequest();
  var queryMethod = queryMethod === undefined ? "GET" : queryMethod;
  var storageKey = storageKey === undefined ? "temp" : storageKey;
  
  httpRequest.open(queryMethod, queryURL + (queryMethod === "GET" ? "?" + $.param(queryParam) : "") , true);
  httpRequest.onload = function(e){
    if (httpRequest.readyState === 4 && httpRequest.status === 200){
      localStorage.setItem(storageKey, httpRequest.responseText);
      if (callback !== undefined && callback !== null)
        callback(JSON.parse(httpRequest.responseText), attribute);
    }
      
    else
      console.error(httpRequest.statusText);
  }
  if (!queryURL.includes("file:///")) httpRequest.send(queryMethod === "GET" ? null : $.param(queryParam));
}

function Initialize(){
  identity().get(function(UUID){
    var localUUID = localStorage.getItem("uuid");
    if (localUUID === null || localUUID === ""){
      document.cookie = "uuid=" + UUID + "; max-age=86400;";
      localStorage.setItem("uuid", UUID);
    } 
  });
  
  var localPoolID = localStorage.getItem("poolid");
  if (localPoolID !== null || localPoolID !== "") {
    
  }
}

function LoadList(dataset, callback, idName, attribute){
  var template = "<li><div class='meta'><p class='track'>TRACK</p><p class='artist'>ARTIST</p></div><div class='control'><img class='vote downvote' src='./assets/downvote.svg' alt='Downvote' data-track='TRACKID' /><img class='vote upvote' src='./assets/upvote.svg' alt='Upvote' data-track='TRACKID' /></div></li>";
}

var tempURL = DatasetURL + "dataset.json";
var tempQueryParam = {};

LoadJSON(tempURL, tempQueryParam, "GET" , "list");

{
  // var playlistid = dataset["playlistid"] === undefined ? "" : dataset["playlistid"];
  // var userid = dataset["userid"] === undefined ? "" : dataset["userid"];
  // localStorage.setItem("playlistid", playlistid);
  // localStorage.setItem("userid", userid);
}