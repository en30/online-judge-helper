(function($){
  $.appendSolveButton().on('click', function() {
    var params ={
      site: 'aoj',
      time_limit: parseInt($('h1.title').next().text().match(/(\d+) sec/)[1], 10),
      id: location.search.match(/id=(.*?)(?:&|$)/m)[1],
      samples: []
    };
    $("h2").each(function(){
      if(!$(this).text().match(/^Sample Input/)) return;
      var $input = $(this).next('pre'),
      $output = $input.next().next('pre');
      params.samples.push({input: $input.text(), output: $output.text()});
    });
    $.solve(params);
  });
})(jQuery)
