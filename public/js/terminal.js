/**
 * Created by satyakb on 7/19/15.
 */

$(function() {
    var conn;

    function launchTerminal() {
        var term = new Terminal({
            cols: 80,
            rows: 40,
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
        launchTerminal();
    } else {
        throw new Error('Your browser does not support WebSockets.');
    }



});
