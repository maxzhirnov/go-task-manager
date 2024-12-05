const API_URL = '/api/tasks';
const USER_API_URL = '/api/users';

export async function handlePublicRequest(url, options = {}) {
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const error = new Error(errorData.error || response.statusText);
            error.status = response.status;
            error.details = errorData;
            throw error;
        }
        if (response.status === 204) return null;
        return await response.json();
    } catch (error) {
        throw error;
    }
}

export async function fetchWithAuth(url, options = {}) {
    try {
        let accessToken = localStorage.getItem("jwt");
        if (!accessToken) {
            window.location.href = "/login";
            return;
        }

        options.headers = {
            ...options.headers,
            "Authorization": `Bearer ${accessToken}`
        };

        let response = await fetch(url, options);

        if (response.status === 401) {
            const newAccessToken = await refreshTokenRequest();
            if (!newAccessToken) {
                window.location.href = "/login";
                return;
            }
            options.headers["Authorization"] = `Bearer ${newAccessToken}`;
            response = await fetch(url, options);
        }
        return response;
    } catch (error) {
        console.error("Request failed:", error);
        throw error;
    }
}

export async function refreshTokenRequest() {
    const refreshToken = localStorage.getItem("refresh_token");
    if (!refreshToken) return null;

    try {
        const response = await fetch("/api/refresh", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ refresh_token: refreshToken })
        });

        if (!response.ok) throw new Error("Failed to refresh token");

        const { access_token } = await response.json();
        localStorage.setItem("jwt", access_token);
        return access_token;
    } catch (error) {
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        return null;
    }
}

export async function handleApiRequest(url, options = {}) {
    try {
        const response = await fetchWithAuth(url, options);
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const error = new Error(errorData.error || response.statusText);
            error.status = response.status;
            error.details = errorData;
            throw error;
        }
        if (response.status === 204) return null;
        return await response.json();
    } catch (error) {
        throw error;
    }
}

async function editTask(taskId, title, description, status) {
    try {
        await handleApiRequest(`${API_URL}/${taskId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, description, status })
        });
        fetchTasks();
    } catch (error) {
        console.error("Failed to edit task:", error);
    }
}

export const api = {
    fetchTasks: async () => {
        const tasks = await handleApiRequest(API_URL);
        // Sort tasks by position in descending order
        return tasks.sort((a, b) => a.position - b.position);
    },
    
    addTask: (title, description, status) => handleApiRequest(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, description, status })
    }),
    
    deleteTask: (id) => handleApiRequest(`${API_URL}/${id}`, { method: "DELETE" }),
    
    updateTaskStatus: async (taskId, newStatus) => {
        const task = await handleApiRequest(`${API_URL}/${taskId}`);
        return handleApiRequest(`${API_URL}/${taskId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                title: task.title,
                description: task.description,
                status: newStatus
            })
        });
    },

    updateTask: async (taskId, title, description, status) => {
        return handleApiRequest(`${API_URL}/${taskId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, description, status })
        });
    },

    updateTaskPositions: async (positions) => {
        return handleApiRequest(`${API_URL}/positions`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(positions)
        });
    },

    getUserStatistics: async () => {
        return handleApiRequest(`${USER_API_URL}/statistics`, {
            method: "GET",
            headers: { "Content-Type": "application/json" }
        });
    },

    updateProfile: async (data) => {
        return handleApiRequest('/api/profile', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
    },

    refreshToken: async () => {
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) {
            console.error('No refresh token found');
            return null;
        }
    
        try {
            console.log('Attempting to refresh token');
            const response = await fetch('/api/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken })
            });
    
            if (!response.ok) {
                throw new Error(`Refresh failed: ${response.statusText}`);
            }
    
            const data = await response.json();
            console.log('Token refresh successful');
            console.log(data);
            
            if (data.access_token) {
                localStorage.setItem('jwt', data.access_token);
                return data.access_token;
            }
            return null;
        } catch (error) {
            console.error('Token refresh failed:', error);
            return null;
        }
    },

    requestPasswordReset: async (data) => {
        return handlePublicRequest('/api/forgot-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
    },

    resetPassword: async (data) => {
        return handlePublicRequest('/api/reset-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
    },
};

