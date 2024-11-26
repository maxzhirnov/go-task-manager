const API_URL = '/api/tasks';

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
        showError(error.message);
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
        showError(error.message);
        localStorage.removeItem("jwt");
        localStorage.removeItem("refresh_token");
        return null;
    }
}

export async function handleApiRequest(url, options = {}) {
    try {
        const response = await fetchWithAuth(url, options);
        if (!response.ok) throw new Error(`API request failed: ${response.statusText}`);
        if (response.status === 204) return null;
        return await response.json();
    } catch (error) {
        showError(error.message);
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
    fetchTasks: () => handleApiRequest(API_URL),
    
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
    }
};

