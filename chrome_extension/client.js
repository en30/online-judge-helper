(function($){
  $.appendSolveButton = function () {
    var $solve = $(document.createElement('div')).text('solve').css({
      cursor: 'pointer',
      position: 'fixed',
      left: 0,
      bottom: 0,
      padding: '1em 2em',
      color: 'white',
      backgroundColor: 'rgba(0, 0, 255, .4)',
      zIndex: 2147483647
    });
    $('body').append($solve);
    return $solve;
  };

  $.solve = function(params) {
    $.post('http://localhost:4567/problem', params);
  };
})(jQuery);
