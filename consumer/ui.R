library(shiny)
library(ggvis)
library(DT)

shinyUI(fluidPage(
   
   titlePanel("Kinesis Data Viewer"),
   
   mainPanel(
      ggvisOutput("lineplot"),
      dataTableOutput("table")
      
   )
))
