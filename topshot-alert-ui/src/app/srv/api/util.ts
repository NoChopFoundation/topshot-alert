import { IPlay } from "./plays-database.service";
import { ISet } from "./sets-database.service";


function shorthandEditionCount(editionCount: number): string {
  if (editionCount % 1000) {
    return `${(editionCount % 1000)}K`
  } else {
    return `${editionCount}`;
  }
}

export function Util_GetSetNameWithEditionInfo_ViaSetId(CompleteSets: ISet[], play: IPlay, setId: number) {
  const setIndex = play.Sets.findIndex(sId => sId === setId);
  return Util_GetSetNameWithEditionInfo_ViaPlaySetIndex(CompleteSets, play, setIndex);
}

export function Util_GetSetNameWithEditionInfo_ViaPlaySetIndex(CompleteSets: ISet[], play: IPlay, setIndexWithinPlaySet: number): string {
  const setId = play.Sets[setIndexWithinPlaySet];
  const editionCount = play.EditionCounts[setIndexWithinPlaySet];
  for (var i = 0; i < CompleteSets.length; i++) { // hash table this
    if (CompleteSets[i].SetId === setId) {
      if (CompleteSets[i].SetSeries > 0) {
        return CompleteSets[i].Name + "(S" + CompleteSets[i].SetSeries + " of " + shorthandEditionCount(editionCount) + ")";
      } else {
        return CompleteSets[i].Name + "(S1 of " + shorthandEditionCount(editionCount) + ")";
      }
    }
  }
  return "UnknownSet[" + setId + "]";
}
