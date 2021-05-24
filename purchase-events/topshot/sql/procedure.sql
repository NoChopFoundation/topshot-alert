DROP procedure IF EXISTS `GetRecentMoments`;

DELIMITER $$
CREATE PROCEDURE `GetRecentMoments`(
	IN pRequestHeight INT,
    IN pMaxBackward INT
)
BEGIN
	DECLARE maxBlock INT Default 0;
    DECLARE fromHeight INT;
    DECLARE toHeight INT;
    
SELECT 
    MAX(BlockHeight)
INTO maxBlock FROM
    moment_events;

    IF maxBlock < pRequestHeight THEN 
		SET fromHeight = maxBlock - pMaxBackward;
        SET toHeight = maxBlock;
	ELSEIF (maxBlock - pRequestHeight) < pMaxBackward THEN
		SET fromHeight = pRequestHeight;
		SET toHeight = maxBlock;
    ELSE
		SET fromHeight = maxBlock - pMaxBackward;
        SET toHeight = maxBlock;
    END  IF;
    
SELECT 
    type,
    MomentId,
    BlockHeight,
    PlayId,
    SerialNumber,
    SetId,
    SellerAddr,
    Price,
    created_at
FROM
    moment_events
WHERE
    BlockHeight <= toHeight
        AND BlockHeight > fromHeight
ORDER BY BlockHeight DESC;
END$$
DELIMITER ;
