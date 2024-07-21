/* Platform Typings */
export interface PlatformInterface {
	get(name: string): void;
}
export interface PlatformType {
	name: string;
	description: string;
	logo: string;
	show: boolean;
	platform: PlatformInterface;
}

/* Plugin Typings */
export interface PluginInterface {
	get(name: string): void;
}
export interface PluginType {
	name: string;
	description: string;
	logo: string;
	show: boolean;
	plugin: PluginInterface;
}
