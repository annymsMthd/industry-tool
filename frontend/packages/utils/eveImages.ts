/**
 * Utility functions for generating EVE Online image URLs
 */

/**
 * Get the URL for an EVE item type icon
 * @param typeId - The EVE type ID
 * @param size - Icon size (32, 64, 128, 256, 512)
 * @returns URL to the item icon
 */
export function getItemIconUrl(typeId: number, size: 32 | 64 | 128 | 256 | 512 = 64): string {
  return `https://images.evetech.net/types/${typeId}/icon?size=${size}`;
}

/**
 * Get the URL for an EVE item type render
 * @param typeId - The EVE type ID
 * @param size - Render size (32, 64, 128, 256, 512)
 * @returns URL to the item render
 */
export function getItemRenderUrl(typeId: number, size: 32 | 64 | 128 | 256 | 512 = 512): string {
  return `https://images.evetech.net/types/${typeId}/render?size=${size}`;
}

/**
 * Get the URL for an EVE character portrait
 * @param characterId - The EVE character ID
 * @param size - Portrait size (32, 64, 128, 256, 512, 1024)
 * @returns URL to the character portrait
 */
export function getCharacterPortraitUrl(characterId: number, size: 32 | 64 | 128 | 256 | 512 | 1024 = 128): string {
  return `https://images.evetech.net/characters/${characterId}/portrait?size=${size}`;
}

/**
 * Get the URL for an EVE corporation logo
 * @param corporationId - The EVE corporation ID
 * @param size - Logo size (32, 64, 128, 256)
 * @returns URL to the corporation logo
 */
export function getCorporationLogoUrl(corporationId: number, size: 32 | 64 | 128 | 256 = 128): string {
  return `https://images.evetech.net/corporations/${corporationId}/logo?size=${size}`;
}

/**
 * Get the URL for an EVE alliance logo
 * @param allianceId - The EVE alliance ID
 * @param size - Logo size (32, 64, 128, 256)
 * @returns URL to the alliance logo
 */
export function getAllianceLogoUrl(allianceId: number, size: 32 | 64 | 128 | 256 = 128): string {
  return `https://images.evetech.net/alliances/${allianceId}/logo?size=${size}`;
}
