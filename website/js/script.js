var BaseURL = "https://linode.shellcode.in/"
var DatasetURL = document.baseURI.substring(0, document.baseURI.lastIndexOf("/") + 1) + "assets/";

function LoadJSON(queryURL, queryParam, queryMethod, storageKey, callback, attribute){
  var httpRequest = new XMLHttpRequest();
  var queryMethod = queryMethod === undefined ? "GET" : queryMethod;
  var storageKey = storageKey === undefined ? "temp" : storageKey;
  
  httpRequest.open(queryMethod, queryURL + (queryMethod === "GET" ? "?" + $.param(queryParam) : "") , true);
  //httpRequest.setRequestHeader("Access-Control-Allow-Origin", "*");
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
}

function JoinPool(){
  
  var inputPoolID = document.getElementById("poolid").value;
  if (inputPoolID !== undefined && inputPoolID !== null && inputPoolID !== ""){
    localStorage.setItem("poolid", inputPoolID);
    var queryResult = LoadJSON(BaseURL + "join_pool", {"poolShortId": inputPoolID, "userId": localStorage.getItem("uuid")}, "POST", "dataset", Initialize, {});
    console.log(queryResult);
  }
}

function QuitPool(){
  localStorage.removeItem("poolid");
  localStorage.removeItem("dataset");
  window.location.reload(true);
}

function LoadList(dataset, idName, callback, attribute){
  var idName = idName === undefined ? "list" : idName;
  var template = "<li><img class='cover' src='COVER' alt='TRACK'/><div class='meta'><p class='track'>TRACK - ALBUM</p><p class='artist'>ARTIST</p></div><div class='control'><span class='count'>DOWNVOTES</span><img class='vote downvote' src='./assets/downvote.svg' alt='Downvote' data-track='TRACKID' /><span class='count'>UPVOTES</span><img class='vote upvote' src='./assets/upvote.svg' alt='Upvote' data-track='TRACKID' /></div></li>";
  
  document.getElementById(idName).innerHTML = "";
  
  for (var i = 0; i < dataset.length; ++i){
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

Initialize();

$("#quit").click(QuitPool);
$("#quit").singletap(QuitPool);