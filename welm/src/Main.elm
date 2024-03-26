module Main exposing (..)

import Browser
import Browser.Navigation as Nav
import Url

import Html exposing (Html, div, button, text, input, span)
import Html.Events exposing (onClick, onInput)
import Debug exposing (toString)
import Html.Attributes exposing (class)
import Html.Attributes exposing (placeholder)
import Html.Attributes exposing (value)
import Html.Attributes exposing (type_)
import Html.Attributes exposing (style)
import Http
import Json.Decode exposing (Decoder, map4, field, int, string)


type alias Model =
  { key : Nav.Key
  , url : Url.Url
  , count : Int
  }


init : () -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
  ({ key = key
  , url = url
  , count = 0 
  }, Cmd.none )


type Msg = LinkClicked Browser.UrlRequest
  | UrlChanged Url.Url
  | Increment 
  | Decrement 
  

update: Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    LinkClicked urlRequest ->
      case urlRequest of
        Browser.Internal url ->
          ( model, Nav.pushUrl model.key (Url.toString url) )

        Browser.External href ->
          ( model, Nav.load href )

    UrlChanged url ->
      ( { model | url = url }
      , Cmd.none
      )

    Increment ->
      ({ model | count = model.count + 1 }, Cmd.none)
    Decrement ->
      ({ model | count = model.count - 1 }, Cmd.none)


view: Model -> Browser.Document Msg
view model =
  {
    title = "Counter"
    , body = [
      div [ class "w-full flex flex-row justify-center items-center" ]
        [
          span [ class "flex flex-row items-center" ] [
            -- button [ class "block w-full rounded-md bg-indigo-600 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600", onClick Decrement ] [ text "-" ]
            viewBtn "-" Decrement
            , div [ class "mx-5" ] [ text ( toString model.count ) ]
            , viewBtn "+" Increment
            -- , button [ onClick Increment ] [ text "+" ]
          ]
        ]
    ] 
  }


viewBtn : String -> Msg -> Html Msg
viewBtn label msg =
  button [ class "block w-full rounded-md bg-indigo-600 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600", onClick msg ] [ text label ]

subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none


main : Program () Model Msg
main = 
  Browser.application 
    { init = init
    , update = update
    , view = view
    , subscriptions = subscriptions 
    , onUrlChange = UrlChanged
    , onUrlRequest = LinkClicked
    }

