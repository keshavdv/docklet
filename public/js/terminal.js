/**
 * Created by satyakb on 7/19/15.
 */

window.onload = function() {
    var socket = io(location.protocol + '//' + window.location.host,{path: '/terminal-ws', query: {id: window.id}});
    socket.on('connect', function() {
        var term = new Terminal({
            cols: 80,
            rows: 40,
            convertEol: true,
            useStyle: true,
            screenKeys: true,
            cursorBlink: false
        });
        term.on('data', function(data) {
            // console.log('[CLIENT]', data);
            //term.write(data);
            socket.emit('input', data);
        });
        term.on('title', function(title) {
            document.title = title;
        });
        term.open(document.body);
        term.write('\x1b[31mWelcome to term.js!\x1b[m\r\n');
        socket.on('output', function(data) {
            // console.log('[SERVER]', data);
            term.write(data);
        });
        socket.on('disconnect', function() {
            term.destroy();
        });
    });
};