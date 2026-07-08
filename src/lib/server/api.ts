import { env } from '$env/dynamic/private';

const base = () => env.API_URL ?? 'http://localhost:3100';

export class ApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.status = status;
	}
}

const readMessage = async (response: Response) => {
	const body = (await response.json().catch(() => null)) as { message?: string } | null;
	return body?.message ?? response.statusText;
};

export const apiGet = async <T>(cookie: string, path: string): Promise<T> => {
	const response = await fetch(`${base()}${path}`, { headers: { cookie } });
	if (!response.ok) throw new ApiError(response.status, await readMessage(response));
	return response.json() as Promise<T>;
};

export const apiSend = async (
	cookie: string,
	method: string,
	path: string,
	body?: unknown
): Promise<{ ok: true } | { ok: false; status: number; message: string }> => {
	const headers: Record<string, string> = { cookie };
	if (body !== undefined) headers['content-type'] = 'application/json';

	const response = await fetch(`${base()}${path}`, {
		method,
		headers,
		body: body === undefined ? undefined : JSON.stringify(body)
	});
	if (response.ok) return { ok: true };
	return { ok: false, status: response.status, message: await readMessage(response) };
};
