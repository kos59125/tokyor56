library(jsonlite)

# default profile を使う

describe_stream <- function(stream) {
   cmd <- sprintf("aws kinesis describe-stream --stream-name %s", stream)
   r <- system(cmd, intern = TRUE)
   fromJSON(paste(r, collapse = ""))
}

get_shared_iterator <- function(stream, shard_id, iterator_type = c("TRIM_HORIZON", "LATEST", "AT_SEQUENCE_NUMBER", "AFTER_SEQUENCE_NUMBER", "AT_TIMESTAMP"), timestamp, sequence_number) {
   iterator_type <- match.arg(iterator_type)
   iterator_type_option <- switch(
      iterator_type,
      "AT_TIMESTAMP" = sprintf("--shard-iterator-type %s --timestamp %.0f", iterator_type, as.numeric(timestamp)),
      "AT_SEQUENCE_NUMBER" = ,
      "AFTER_SEQUENCE_NUMBER" = sprintf("--shard-iterator-type %s --starting-sequence-number %s", iterator_type, sequence_number),
      sprintf("--shard-iterator-type %s", iterator_type)
   )
   
   cmd <- sprintf("aws kinesis get-shard-iterator --stream-name %s --shard-id %s %s", stream, shard_id, iterator_type_option)
   message(cmd)
   r <- system(cmd, intern = TRUE)
   fromJSON(paste(r, collapse = ""))
}

get_records <- function(iterator, limit) {
   extra <- ""
   if (!missing(limit)) {
      extra <- sprintf(" --limit %d", limit)
   }
   
   cmd <- sprintf("aws kinesis get-records --shard-iterator %s%s", shQuote(iterator), extra)
   message(cmd)
   r <- system(cmd, intern = TRUE)
   fromJSON(paste(r, collapse = ""))
}
