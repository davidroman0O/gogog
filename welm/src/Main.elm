module Main exposing (..)

import Browser
import Html exposing (Html, div, button, text, input)
import Html.Events exposing (onClick, onInput)
import Debug exposing (toString)
import Html.Attributes exposing (class)
import Html.Attributes exposing (placeholder)
import Html.Attributes exposing (value)
import Html.Attributes exposing (type_)
import Html.Attributes exposing (style)
import Http
import Json.Decode exposing (Decoder, map4, field, int, string)


initModel: () -> (Model, Cmd Msg)
initModel _ = 
  ()

type Msg = 
  Increment 
  | Decrement 
  

update: Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    Increment ->
      ({ model | count = model.count + 1 }, Cmd.none)
    Decrement ->
      ({ model | count = model.count - 1 }, Cmd.none)
      

view model =
  div [ class "w-full flex flex-row justify-center items-center bg-red-500" ]
    [
      button [ class "bg-red-500 rounded", onClick Decrement ] [ text "-" ]
      , div [] [ text ( toString model.count ) ]
      , button [ onClick Increment ] [ text "+" ]
    ]


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none


main = 
  Browser.application 
    { init = initModel
    , update = update
    , view = view
    , subscriptions = subscriptions }

