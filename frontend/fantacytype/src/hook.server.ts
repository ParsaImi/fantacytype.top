import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    // Add any authentication checking logic here
    return resolve(event);
};
