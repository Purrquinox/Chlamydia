/* Global Typings */
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
