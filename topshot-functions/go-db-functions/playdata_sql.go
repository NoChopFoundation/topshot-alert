package playdata

func GetSetSelect(opts *GetPlayDataHTTPInit) string {
	return "SELECT SetId, Name, SetSeries FROM " + DBTABLE_sets
}

func GetSetSelectWithPlayerData(opts *GetPlayDataHTTPInit) string {
	return "SELECT " + DBTABLE_plays_in_sets + ".PlayId, " + DBTABLE_plays + ".FullName " +
		"FROM " + DBTABLE_plays_in_sets + " INNER JOIN " + DBTABLE_plays + " ON " +
		DBTABLE_plays_in_sets + ".PlayId = " + DBTABLE_plays + ".PlayId"
}

func GetPlaySelect(opts *GetPlayDataHTTPInit) string {
	return "SELECT PlayId, NbaSeason, TeamAtMomentNBAID, PlayCategory, JerseyNumber, PlayerPosition, " +
		"DateOfMoment, PlayType, FullName, PrimaryPosition, TeamAtMoment FROM " + DBTABLE_plays
}

func GetPlaysInSets(opts *GetPlayDataHTTPInit) string {
	return "SELECT PlayId, SetId, EditionCount FROM " + DBTABLE_plays_in_sets
}

func GetRecentCollectorStatus(opts *GetPlayDataHTTPInit) string {
	return "SELECT CollectorId, State, UpdatesInInterval, BlockHeight, CreatedAt from " + DBTABLE_moment_events_collectors + " WHERE " +
		"`CreatedAt` >= date_sub(NOW(), INTERVAL ? SECOND) order by CreatedAt desc"
}

func GetRecentByPlaySet(opts *GetPlayDataHTTPInit) string {
	return "select type, MomentId, BlockHeight, PlayId, SerialNumber, SetId, SellerAddr, Price, Created_At from " + DBTABLE_moment_events +
		" where PlayId=? and SetId=? and " +
		"BlockHeight > ((select MAX(BlockHeight) from " + DBTABLE_moment_events + ") - ?) order by BlockHeight desc;"
}

// Returns type, MomentId, BlockHeight, PlayId, SerialNumber, SetId, SellerAddr
func GetRecentMoments(opts *GetPlayDataHTTPInit) string {
	return "CALL GetRecentMoments(?, ?)"
}
