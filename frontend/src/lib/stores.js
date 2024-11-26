import { writable } from 'svelte/store';

export const tasks = writable([]);
export const errorMessage = writable('');
export const user = writable(null);

export function showError(message) {
    errorMessage.set(message);
    setTimeout(() => errorMessage.set(''), 3000);
}