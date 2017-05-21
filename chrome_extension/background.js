var graphGenerator = function(directed) {
  return function(obj, tab){
    // object.selectedText does not include new line characters
    chrome.tabs.executeScript(tab.id, {
      code: "window.getSelection().toString()"
    }, function(selection){
      var xhr = new XMLHttpRequest();
      xhr.open('POST', 'http://localhost:4567/graph');
      xhr.responseType = 'blob';
      xhr.onload = function() {
        window.open(URL.createObjectURL(this.response));
      };
      xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
      xhr.send($.param({ directed: directed, adjacentList: selection[0]}));
    });
  };
};

var rootId = chrome.contextMenus.create({
  title: 'Online Judge Helper',
  contexts: ['selection']
});

chrome.contextMenus.create({
  title: 'Undirected Graph',
  contexts: ['selection'],
  parentId: rootId,
  onclick: graphGenerator(false)
});

chrome.contextMenus.create({
  title: 'Directed Graph',
  contexts: ['selection'],
  parentId: rootId,
  onclick: graphGenerator(true)
});
