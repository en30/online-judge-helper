(function($){
  $.appendSolveButton().on('click', function() {
    var params ={
      site: 'poj',
      time_limit: Number($('.plm').text().match(/(\d+)MS/)[1])/1000,
      id: location.search.match(/id=(.*?)(?:&|$)/m)[1],
      samples: []
    },
        ios = $('.sio').map(function(){ return $(this).text(); });
    params.samples.push({input: ios[0], output: ios[1]});
    $.solve(params);
  });
})(jQuery);
