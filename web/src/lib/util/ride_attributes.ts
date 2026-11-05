// Bike-only fork: labels for the ride attributes that replace "difficulty".
// Labels are German on purpose — this instance has a single, German-speaking user.

export type RouteType = "loop" | "traversal";
export type RideProfile =
    | "road-low-traffic"
    | "ffm-by-bicycle"
    | "trekking"
    | "quaelnix-gravel"
    | "m11n-gravel"
    | "other";

export const ROUTE_TYPES: { value: RouteType; text: string; icon: string }[] = [
    { value: "loop", text: "Rundkurs", icon: "rotate-right" },
    { value: "traversal", text: "A → B", icon: "route" },
];

export const RIDE_PROFILES: { value: RideProfile; text: string }[] = [
    { value: "road-low-traffic", text: "Rennrad (sehr wenig Verkehr)" },
    { value: "ffm-by-bicycle", text: "FFM by Bicycle" },
    { value: "trekking", text: "Trekking" },
    { value: "quaelnix-gravel", text: "Gravel „Quaelnix“" },
    { value: "m11n-gravel", text: "Gravel „m11n mehr Offroad Beta“" },
    { value: "other", text: "Anderes Profil" },
];

export function routeTypeLabel(value?: string | null): string {
    return ROUTE_TYPES.find((r) => r.value === value)?.text ?? "Streckentyp ?";
}

export function routeTypeIcon(value?: string | null): string {
    return ROUTE_TYPES.find((r) => r.value === value)?.icon ?? "route";
}

export function rideProfileLabel(value?: string | null): string {
    return RIDE_PROFILES.find((p) => p.value === value)?.text ?? "Profil ?";
}
