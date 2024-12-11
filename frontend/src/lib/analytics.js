import mixpanel from 'mixpanel-browser';

// Initialize Mixpanel
mixpanel.init('9ecff9b2013eeaeda0a3c0635fa9e084', {
    debug: false, // import.meta.env.DEV,
    track_pageview: true,
    persistence: 'localStorage',
    ignore_dnt: true,
});

let currentUser = null;

export const Analytics = {
    track: (event_name, properties = {}) => {
        // Always include user properties in every event if user exists
        const userProperties = currentUser ? {
            distinct_id: currentUser.id,
            user_id: currentUser.id,
            email: currentUser.email,
            username: currentUser.username,
        } : {};

        mixpanel.track(event_name, {
            ...userProperties,
            ...properties,
        });
    },
    
    identify: (userId) => {
        mixpanel.identify(userId);
    },
    
    setUserProperties: (properties) => {
        mixpanel.people.set(properties);
    },
    
    pageView: (page) => {
        Analytics.track('Page View', { page });
    },

    trackUser: (user) => {
        if (!user) return;

        currentUser = user; // Store user info
        mixpanel.identify(user.id);
        mixpanel.people.set({
            $email: user.email,
            $name: user.username,
            $last_login: new Date().toISOString(),
        });

        // Analytics.track('User Logged In', {
        //     loginTime: new Date().toISOString()
        // });
    },

    clearUser: () => {
        currentUser = null;
        mixpanel.reset();
    }
};