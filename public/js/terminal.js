/**
 * Created by satyakb on 7/19/15.
 */

function getCharSize() {
    var $span = $("<span>", {text: "qwertyuiopasdfghjklzxcvbnm"});
    $('#terminal').append($span);
    var size = {
        width: $span.outerWidth()/26
        , height: $span.outerHeight()
    };
    $span.remove();
    return size;
}
function getwindowSize() {
    var e = window,
        a = 'inner';
    if (!('innerWidth' in window )) {
        a = 'client';
        e = document.documentElement || document.body;
    }
    return {width: e[a + 'Width'], height: e[a + 'Height']};
}

function textSize() {
    var charSize = {width: 6.7, height: 17} ;
        //getCharSize();
    var windowSize = getwindowSize();
    return {
        x: Math.floor(windowSize.width / charSize.width)
        , y: Math.floor(windowSize.height / charSize.height)
    };
}

$(function() {
    var conn;

    function launchTerminal() {
        var rc = textSize();
        console.log(rc);
        var term = new Terminal({
            rows: rc.y,
            cols: rc.x,
            convertEol: true,
            useStyle: true,
            screenKeys: true,
            cursorBlink: false
        });

        term.on('data', function(data) {
            conn.send(data);
        });
        term.on('title', function(title) {
            document.title = title;
        });
        term.open(document.body);
        term.write('\x1b[31mWelcome to term.js!\x1b[m\r\n');
        conn.onmessage = function(evt) {
            term.write(evt.data);
        }
        conn.onclose = function(evt) {
            term.destroy();
        }
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + window.location.host + "/terminal-ws?id=" + window.id);
        console.log('launching...');
        launchTerminal();
        console.log('launched');
    } else {
        throw new Error('Your browser does not support WebSockets.');
    }



});
