import type { Group } from '../types/worldcup';

let groupsCache: Group[] = [];

export function setGroups(data: Group[]) {
  groupsCache = data;
}

export function getGroups(): Group[] {
  return groupsCache;
}

export function getGroupByName(name: string): Group | undefined {
  return groupsCache.find((g) => g.name === name);
}
