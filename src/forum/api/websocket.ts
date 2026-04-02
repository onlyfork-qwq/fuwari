import { forumEnv } from "@/forum/stores/env";
import { get } from "svelte/store";

export interface WebSocketMessage {
	type: string;
	payload?: Record<string, unknown>;
}

export interface NewCommentPayload {
	postId: string;
	comment: {
		content: string;
		author_name: string;
		author_id: number;
		parent_id: string | null;
		created_at: string;
	};
}

export interface PostUpdatedPayload {
	postId: string;
	title: string;
	content: string;
	category_id: number;
	updated_at: string;
}

export interface SubscribedPayload {
	postId: string;
}

export interface ConnectedPayload {
	message: string;
}

export type WebSocketEventHandler = (payload: Record<string, unknown>) => void;

export class ForumWebSocket {
	private ws: WebSocket | null = null;
	private baseUrl: string;
	private reconnectAttempts = 0;
	private maxReconnectAttempts = 5;
	private reconnectDelay = 1000;
	private pingInterval: ReturnType<typeof setInterval> | null = null;
	private eventHandlers: Map<string, WebSocketEventHandler[]> = new Map();
	private currentPostId: string | null = null;
	private isConnecting = false;

	constructor() {
		this.baseUrl = this.buildWebSocketUrl();
	}

	private buildWebSocketUrl(): string {
		const baseUrl = get(forumEnv.baseUrl);
		const url = new URL(baseUrl);
		const wsProtocol = url.protocol === "https:" ? "wss:" : "ws:";
		return `${wsProtocol}//${url.host}`;
	}

	connect(postId?: string): void {
		if (
			this.isConnecting ||
			(this.ws && this.ws.readyState === WebSocket.OPEN)
		) {
			return;
		}

		this.isConnecting = true;
		this.currentPostId = postId || null;

		const url = postId
			? `${this.baseUrl}/api/ws?postId=${encodeURIComponent(postId)}`
			: `${this.baseUrl}/api/ws`;

		try {
			this.ws = new WebSocket(url);

			this.ws.onopen = () => {
				console.log("WebSocket connected");
				this.isConnecting = false;
				this.reconnectAttempts = 0;
				this.startPingInterval();
				this.emit("connected", { message: "WebSocket connected" });
			};

			this.ws.onmessage = (event) => {
				try {
					const data = JSON.parse(event.data) as WebSocketMessage;
					this.handleMessage(data);
				} catch (error) {
					console.error("Failed to parse WebSocket message:", error);
				}
			};

			this.ws.onclose = () => {
				console.log("WebSocket disconnected");
				this.isConnecting = false;
				this.stopPingInterval();
				this.emit("disconnected", {});
				this.scheduleReconnect();
			};

			this.ws.onerror = (error) => {
				console.error("WebSocket error:", error);
				this.isConnecting = false;
				this.emit("error", { error: error.type });
			};
		} catch (error) {
			console.error("Failed to create WebSocket connection:", error);
			this.isConnecting = false;
			this.scheduleReconnect();
		}
	}

	private handleMessage(data: WebSocketMessage): void {
		switch (data.type) {
			case "connected":
				this.emit("connected", data.payload || {});
				break;
			case "subscribed":
				this.emit("subscribed", data.payload || {});
				break;
			case "new_comment":
				this.emit("new_comment", data.payload || {});
				break;
			case "post_updated":
				this.emit("post_updated", data.payload || {});
				break;
			case "pong":
				this.emit("pong", {});
				break;
			default:
				console.log("Unknown WebSocket message type:", data.type);
		}
	}

	subscribePost(postId: string): void {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(
				JSON.stringify({
					type: "subscribe_post",
					payload: { postId },
				}),
			);
			this.currentPostId = postId;
		}
	}

	unsubscribePost(): void {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(
				JSON.stringify({
					type: "unsubscribe_post",
					payload: {},
				}),
			);
			this.currentPostId = null;
		}
	}

	ping(): void {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify({ type: "ping" }));
		}
	}

	private startPingInterval(): void {
		this.stopPingInterval();
		this.pingInterval = setInterval(() => {
			this.ping();
		}, 30000);
	}

	private stopPingInterval(): void {
		if (this.pingInterval) {
			clearInterval(this.pingInterval);
			this.pingInterval = null;
		}
	}

	private scheduleReconnect(): void {
		if (this.reconnectAttempts < this.maxReconnectAttempts) {
			this.reconnectAttempts++;
			const delay =
				this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
			console.log(`Reconnecting in ${delay}ms...`);
			setTimeout(() => {
				this.connect(this.currentPostId || undefined);
			}, delay);
		} else {
			console.error("Max reconnection attempts reached");
		}
	}

	disconnect(): void {
		this.stopPingInterval();
		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
		this.isConnecting = false;
		this.reconnectAttempts = 0;
	}

	on(event: string, handler: WebSocketEventHandler): void {
		if (!this.eventHandlers.has(event)) {
			this.eventHandlers.set(event, []);
		}
		this.eventHandlers.get(event)!.push(handler);
	}

	off(event: string, handler: WebSocketEventHandler): void {
		const handlers = this.eventHandlers.get(event);
		if (handlers) {
			const index = handlers.indexOf(handler);
			if (index !== -1) {
				handlers.splice(index, 1);
			}
		}
	}

	private emit(event: string, payload: Record<string, unknown>): void {
		const handlers = this.eventHandlers.get(event);
		if (handlers) {
			handlers.forEach((handler) => {
				try {
					handler(payload);
				} catch (error) {
					console.error(
						`Error in WebSocket event handler for ${event}:`,
						error,
					);
				}
			});
		}
	}

	isConnected(): boolean {
		return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
	}

	getCurrentPostId(): string | null {
		return this.currentPostId;
	}
}

let forumWebSocketInstance: ForumWebSocket | null = null;

export function getForumWebSocket(): ForumWebSocket {
	if (!forumWebSocketInstance) {
		forumWebSocketInstance = new ForumWebSocket();
	}
	return forumWebSocketInstance;
}

export function disconnectForumWebSocket(): void {
	if (forumWebSocketInstance) {
		forumWebSocketInstance.disconnect();
		forumWebSocketInstance = null;
	}
}
