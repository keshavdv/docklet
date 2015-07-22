/**
 * Created by satyakb on 7/19/15.
 */

$(function() {
    var conn;
    //var msg = $("#msg");
    //var log = $("#log");
    //function appendLog(msg) {
    //    var d = log[0]
    //    var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
    //    msg.appendTo(log)
    //    if (doScroll) {
    //        d.scrollTop = d.scrollHeight - d.clientHeight;
    //    }
    //}
    //$("#form").submit(function() {
    //    if (!conn) {
    //        return false;
    //    }
    //    if (!msg.val()) {
    //        return false;
    //    }
    //    conn.send(msg.val());
    //    msg.val("");
    //    return false;
    //});

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
            term.write(evt)
        }
        conn.onclose = function(evt) {
            term.destroy();
        }
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + window.location.host + "/terminal-ws");
        launchTerminal();
        //conn.onclose = function(evt) {
        //    appendLog($("<div><b>Connection closed.</b></div>"))
        //}
        //conn.onmessage = function(evt) {
        //    appendLog($("<div/>").text(evt.data))
        //}
    } else {
        //appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
        throw new Error('Your browser does not support WebSockets.');
    }



});


//window.onload = function() {
//    var socket = io(location.protocol + '//' + window.location.host,{path: '/terminal-ws?id=93a7eaccb953'});
//    socket.on('connect', function() {
//        var term = new Terminal({
//            cols: 80,
//            rows: 40,
//            convertEol: true,
//            useStyle: true,
//            screenKeys: true,
//            cursorBlink: false
//        });
//        term.on('data', function(data) {
//            // console.log('[CLIENT]', data);
//            //term.write(data);
//            socket.emit('input', data);
//        });
//        term.on('title', function(title) {
//            document.title = title;
//        });
//        term.open(document.body);
//        term.write('\x1b[31mWelcome to term.js!\x1b[m\r\n');
//        socket.on('output', function(data) {
//            // console.log('[SERVER]', data);
//            term.write(data);
//        });
//        socket.on('disconnect', function() {
//            term.destroy();
//        });
//    });
//};