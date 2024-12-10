import mixpanel from 'mixpanel-browser';

// Initialize Mixpanel
mixpanel.init('9ecff9b2013eeaeda0a3c0635fa9e084', {
    debug: import.meta.env.DEV, // Enable debug in development
    track_pageview: true,
    persistence: 'localStorage'
});

// Helper functions
export const Analytics = {
    track: (event_name, properties = {}) => {
        mixpanel.track(event_name, properties);
    },
    
    identify: (userId) => {
        mixpanel.identify(userId);
    },
    
    setUserProperties: (properties) => {
        mixpanel.people.set(properties);
    },
    
    pageView: (page) => {
        mixpanel.track('Page View', { page });
    }
};