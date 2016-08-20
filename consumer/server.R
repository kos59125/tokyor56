library(shiny)
library(parallel)
library(jsonlite)
library(dplyr)
library(ggvis)
library(DT)

source("kinesis.R", encoding = "UTF-8")

shinyServer(function(input, output, session) {
   
   stream_name <- "accesslog"
   
   shards <- NULL
   iterators <- NULL
   access_log <- NULL
   
   update_stream <- function() {
      stream <- describe_stream(stream_name)
      shards <<- stream$StreamDescription$Shards
      
      it <- sapply(seq_along(shards$ShardId), function(index) {
         id <- shards$ShardId[index]
         get_shared_iterator(stream_name, id)$ShardIterator
      })
      names(it) <- shards$ShardId
      iterators <<- it
   }
   
   fetch_records <- function() {
      records <- mclapply(iterators, function(it) get_records(it, 100))
      iterators <<- sapply(records, function(shard) shard$NextShardIterator)
      
      json <- unlist(lapply(records, function(x) {
         lapply(x$Records$Data, function(j) {
            rawToChar(base64_dec(j))
         })
      }))
      
      if (!is.null(json)) {  # empty records
         newdata <- bind_rows(mclapply(json, function(j) as.data.frame(fromJSON(j), stringsAsFactors = FALSE))) %>% 
            mutate(timestamp = as.POSIXct(sub("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}).*", "\\1", timestamp), tz = "UTC", format = "%Y-%m-%dT%H:%M:%S"))
         access_log <<- bind_rows(access_log, newdata) %>% 
            arrange(timestamp) %>% 
            tail(500)
      }
   }
   
   initialize <- function() {
      update_stream()
      fetch_records()
   }
   initialize()
   
   df <- reactive({
      # 10 秒ごとに更新
      invalidateLater(10000, session)
      fetch_records()
      
      access_log %>% 
         group_by(timestamp = as.POSIXct(10 * (as.numeric(timestamp) %/% 10), origin = "1970-01-01", tz = "UTC")) %>% 
         summarize(count = n())
   })
   
   output$table <- renderDataTable({
      df()
   })
   
   df %>% 
      ggvis(~timestamp, ~count) %>% 
      layer_lines() %>% 
      bind_shiny("lineplot", "lineplot_ui")
   
})
