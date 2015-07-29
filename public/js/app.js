var React = require("react"),
    Fluxxor = require("fluxxor");

window.React = React;

var constants = {
    CREATE_CONTAINER: 'CREATE_CONTAINER',
    SET_STATUS: 'SET_STATUS',
    LOAD_CONTAINER: 'LOAD_CONTAINER'
};

var AppStore = Fluxxor.createStore({
    initialize: function() {
        this.currentImageName = '';
        this.currentContainerID = '';
        this.currentJobID = '';
        this.isPending = false; // if true show loading box
        this.isContainerReady = false;
        this.containerStatus = 'pending';

        this.bindActions(
            constants.CREATE_CONTAINER, this.onCreateContainer,
            constants.LOAD_CONTAINER, this.onLoadContainer,
            constants.SET_STATUS, this.onSetStatus
        )
    },

    onCreateContainer: function(data) {
        this.currentImageName = data.imageName;
        this.currentJobID = data.jobID;
        this.isPending = true;
        this.emit('change');
    },

    onLoadContainer: function(data) {
        this.currentContainerID = data.containerID;
        this.isPending = false;
        this.isContainerReady = true;
        this.emit('change');
    },

    onSetStatus: function(status) {
        this.containerStatus = status;
        this.emit('change');
    },

    getState: function() {
        return {
            currentImageName: this.currentImageName,
            currentContainerID: this.currentContainerID
        };
    }
});


var stores = {AppStore: new AppStore()};
var actions = {
    createImage: function(imageName) {
        var self = this;
        $.ajax({
            url: '/create?image=' + imageName,
            success: function(data) {
                console.log(data);
                self.dispatch(constants.CREATE_CONTAINER, {
                    imageName: imageName,
                    jobID: data.job
                });
            },
            dataType: 'json'
        });
    },
    poll: function(jobID) {
        var self = this;
        var count = 0;
        (function poll() {
            setTimeout(function() {
                $.ajax({
                    url: "/status?job_id=" + jobID,
                    success: function(data) {
                        console.log(data);
                        self.dispatch(constants.SET_STATUS, data.status);
                    },
                    dataType: "json",
                    complete: function(data) {
                        if (data.responseJSON.status === 'success') {
                            start.bind(self)({
                                containerID: data.responseJSON.data.id
                            });
                        } else if (data.responseJSON.status === 'failed') {
                            console.log('Sorry something went wrong...');
                            console.log('Refresh to try again');
                        } else {
                            count++;
                            poll();
                        }
                    }
                });
            }, 3000);
        })();
    }
};

function start(payload) {
    var self = this;
    $.ajax({
        url: "/start?id=" + payload.containerID,
        success: function(data) {
            self.dispatch(constants.LOAD_CONTAINER, payload);
        }
    })
}

var flux = new Fluxxor.Flux(stores, actions);

window.flux = flux;

flux.on("dispatch", function(type, payload) {
    if (console && console.log) {
        console.log("[Dispatch]", type, payload);
    }
});

var FluxMixin = Fluxxor.FluxMixin(React),
    StoreWatchMixin = Fluxxor.StoreWatchMixin;

var PendingBox = React.createClass({
    displayName: "PendingBox",

    mixins: [FluxMixin, StoreWatchMixin('AppStore')],

    getInitialState: function() {
        return {};
    },

    getStateFromFlux: function() {
        return this.getFlux().store('AppStore').getState();
    },

    componentDidMount: function() {
        var jobID = this.getFlux().store('AppStore').currentJobID;
        this.getFlux().actions.poll(jobID);
    },

    render: function() {
        var status = this.getFlux().store('AppStore').containerStatus;
        return (
            <div className="pending-box">
                <span>Status: {status}</span>
            </div>
        );
    }
});

var Form = React.createClass({
    displayName: "Form",

    mixins: [FluxMixin],

    getInitialState: function() {
        return {
            imageName: ''
        }
    },

    onLaunchChange: function(e) {
        this.setState({imageName: e.target.value});
    },

    onSubmit: function(e) {
        e.preventDefault();
        var imageName = this.state.imageName;
        this.getFlux().actions.createImage(imageName);
    },

    render: function() {
        return (
            <form onSubmit={this.onSubmit}>
                <input className="launch" type="text" onChange={this.onLaunchChange} value={this.state.imageName} placeholder="Enter Docker image name from registry..." />
            </form>
        );
    }
});

var Terminal = React.createClass({
    displayName: 'Terminal',

    render: function() {
        var src = '/terminal?id=' + this.props.containerID;
        return (
            <iframe id='terminal-frame' src={src}></iframe>
        );
    }
});

var Application = React.createClass({displayName: "Application",
    mixins: [FluxMixin, StoreWatchMixin('AppStore')],

    getInitialState: function() {
        return {};
    },

    getStateFromFlux: function() {
        return this.getFlux().store('AppStore').getState();
    },
    render: function() {
        var isPending = this.getFlux().store('AppStore').isPending;
        var isReady = this.getFlux().store('AppStore').isContainerReady;
        var containerID = this.getFlux().store('AppStore').currentContainerID;
        var content = isPending ? <PendingBox /> :
            isReady ? <Terminal containerID={containerID} /> : <Form />;
        return (
            <div className="container">
                <div className="wrapper">
                    <div className="row logo-wrapper">
                        <img className="logo" src="images/box.png"></img>
                        <span>Docklet</span>
                    </div>
                    {content}
                </div>
            </div>
        );
    },
});


React.render(React.createElement(Application, {flux: flux}), document.getElementById("app"));