import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';

export function load() {
    if (browser && !localStorage.getItem('jwt')) {
        throw redirect(303, '/');
    }
}