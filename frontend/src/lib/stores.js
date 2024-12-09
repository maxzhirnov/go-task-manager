import { writable } from 'svelte/store';

export const tasks = writable([]);
export const errorMessage = writable('');
export const successMessage = writable('');
export const user = writable(null);

export function showError(message) {
    errorMessage.set(message);
    setTimeout(() => errorMessage.set(''), 3000);
}

export function showSuccess(message) {
    successMessage.set(message);
    setTimeout(() => successMessage.set(''), 3000); // Hide after 3 seconds
}

export function parseJWT(token) {
    if (!token) return null;
    
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        return JSON.parse(jsonPayload);
    } catch (error) {
        console.error('Failed to parse JWT:', error);
        return null;
    }
}

export function initializeUser() {
    const token = localStorage.getItem("jwt");
    if (!token) {
        window.location.href = "/login";
        return;
    }

    const userData = parseJWT(token);
    if (userData) {
        user.set({
            id: userData.user_id,
            username: userData.username,
            email: userData.email
        });
    }
}

// Function to update user data and token
export async function updateUserAndToken() {
    const token = localStorage.getItem("jwt");
    if (token) {
        const userData = parseJWT(token);
        if (userData) {
            user.set({
                id: userData.user_id,
                username: userData.username,
                email: userData.email,
            });
        }
    }
}