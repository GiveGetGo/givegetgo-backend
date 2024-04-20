import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import { Provider } from 'react-redux';
import store from './store';
import LoginScreen from './screens/LoginScreen';
import ForgotPasswordScreen from './screens/ForgotPasswordScreen';
import SignUpScreen from './screens/SignUpScreen'; 
import CheckEmailScreen from './screens/CheckEmailScreen'; 
import ConfirmationScreen from './screens/ConfirmationScreen'; 
import MainScreen from './screens/MainScreen'; 
import GiveOutContactScreen from './screens/GiveOutContactScreen'; 

export type RootStackParamList = {
  LoginScreen: undefined;
  ForgotPasswordScreen: undefined;
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  ConfirmationScreen: undefined;
  MainScreen: undefined;
  GiveOutContactScreen: undefined;
};

const MainStack = createStackNavigator<RootStackParamList>();

const App: React.FC = () => {
  return (
    <Provider store={store}> 
      <NavigationContainer>
        <MainStack.Navigator initialRouteName="LoginScreen">
          <MainStack.Screen name="LoginScreen" component={LoginScreen} options={{ headerShown: false }} />
          <MainStack.Screen name="ForgotPasswordScreen" component={ForgotPasswordScreen} options={{ title: 'Forgot Password', headerShown: false }} />
          <MainStack.Screen name="SignUpScreen" component={SignUpScreen} options={{ title: 'Sign Up', headerShown: false }} />
          <MainStack.Screen name="CheckEmailScreen" component={CheckEmailScreen} options={{ title: 'Check Email', headerShown: false }} />
          <MainStack.Screen name="ConfirmationScreen" component={ConfirmationScreen} options={{ title: 'Confirm', headerShown: false }} />
          <MainStack.Screen name="MainScreen" component={MainScreen} options={{ title: 'Main', headerShown: false }} />
          <MainStack.Screen name="GiveOutContactScreen" component={GiveOutContactScreen} options={{ title: 'GiveOutContact', headerShown: false }} />
          {/* You can add more screens to the navigator as needed */}
        </MainStack.Navigator>
      </NavigationContainer>
    </Provider>
  );
};

export default App;


// Comments in return function will lead to error: "text should be in..."
// View/SafeAreaView is needed when there are more than one components
// tips:　Navigation 放stack's name not component
// for going back to previous page:   const navigation = useNavigation(); <Appbar.BackAction onPress={() => navigation.goBack()} />; (check SeeRequestScreen)
// justifyContent: 'center', // Center contents vertically // alignItems: 'center', // Center contents horizontally
// following up, those positioned "absolute" will not be counted in when using justifyContent or alignItems

// Main tasks:
// Separate MainScreen's four cards so that CSS will be easier to build
// SeeRequestScreen, PostDetailsScreen, NotificationStackProfileScreen: load from json, not redux? (see next line)
// substitute the selected one to './profile_icon.jpg' in all screens (just search './profile_icon.jpg' to find all) (redux) (see SettingsScreen for require() mapping)
// screens裡的圖片改到assests裡
// fill in xxx and some other "Jimmy Ho" for notification stack, home stack pages (see confirmationScreen) (start from the opened screens, then work on the rest while modifying css)
// add userId to redux and screens (got from registration; required in all interactive calls)
// Notifications: find way to hide double-headers (add { headerShown: false } like line 31)
// nativeWind + gpt (no more, just do css)

// if have time: 
// from profile (and some others), make each post's detailed page
// ProfileScreen bio, reply needs word limit
// animation among notification pages? (react-native-reanimated transition; ask claude) (or simply use modal; https://reactnative.dev/docs/modal)
// font?

// CSS status:
// ProfileScreen: (DONE) container, headContainer, header could be replicated to other pages with backspace; card settings
// SeeRequestScreen: (DONE) card in the MIDDLE; avatar; "Take" button
// RequestSucceedScreen: (DONE) header with no backspace (parameters in the 4 sets would exactly match the header w bs); mind each "padding", "marginBottom", "marginTop" values
// GiveOutContactScreen:  (DONE)
// (above are all the notification->request's sub pages)



// when working on css, make sure that in all cases the expected layout works. might have to use a doc

// good card settigns in profileScreen; should unite each component's setting in the end

//profile: (remember that two profileScreens should have different sources) 




// run npm, push code



