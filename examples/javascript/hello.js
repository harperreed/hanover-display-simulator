var FlipDot = require('flipdot-display');
var dateFormat = require('dateformat');

var flippy = new FlipDot('/dev/pts/7',1,16,96);

flippy.on("error", function(err) {
  console.log(err);
});

flippy.once("open", function() {
    str = "Hello There"
    flippy.fill(0xFF);
    font = 'Banner'
    offset = [0,0]
    invert = false
    flippy.writeText(str, {font: font}, offset, invert);

});
