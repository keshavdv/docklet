var React = require("react"),
    Fluxxor = require("fluxxor");

window.React = React;

var stores = {};
var actions = {};

var flux = new Fluxxor.Flux(stores, actions);

window.flux = flux;

flux.on("dispatch", function(type, payload) {
    if (console && console.log) {
        console.log("[Dispatch]", type, payload);
    }
});

var FluxMixin = Fluxxor.FluxMixin(React),
    StoreWatchMixin = Fluxxor.StoreWatchMixin;

var Application = React.createClass({displayName: "Application",
    mixins: [FluxMixin, StoreWatchMixin("TodoStore")],

    getInitialState: function() {
        return {};
    },
    getStateFromFlux: function() {
        return;
    },
    render: function() {
        var todos = this.state.todos;
        return (
            <div className="container">
                <div className="wrapper">
                    <div className="row logo-wrapper">
                        <img className="logo" src="images/box.png"></img>
                        <span>Docklet</span>
                    </div>
                    <input className="launch" type="text" placeholder="Enter Docker image name from registry..." />
                </div>
            </div>
        );
    },
});


React.render(React.createElement(Application, {flux: flux}), document.getElementById("app"));