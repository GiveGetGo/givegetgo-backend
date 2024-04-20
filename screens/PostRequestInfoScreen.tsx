import React, { useState } from 'react';
import { View, StyleSheet, SafeAreaView, Keyboard, TouchableWithoutFeedback  } from 'react-native';
import { Button, Text, Card, TextInput, Appbar } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';

type RootStackParamList = {
  HomeScreen: undefined;
  PostRequestInfoScreen: undefined;
  PostRequestSucceedScreen: undefined;
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

const PostRequestInfoScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {
  const use_navigation = useNavigation(); //for Appbar.BackAction

  const [request, setRequest] = React.useState('');

  const submitRequest = () => {
    console.log('Request sent:', request);
    // Call API to store the 'request' text 
    navigation.navigate('PostRequestSucceedScreen');
  };

  return (
    <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
        <SafeAreaView  style={styles.container}> 
            <View style={styles.headerContainer}>
                <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
                <Text style={styles.header}>GiveGetGo</Text>
                <View style={styles.backActionPlaceholder} />
            </View>
            <View style={styles.request_container}>
                <Card style={styles.card}>
                    <Card.Title title="Request Info" titleStyle={styles.cardTitle}/>
                    <Card.Content>
                    <TextInput
                        label="Please provide details here..."
                        value={request}
                        onChangeText={setRequest}
                        multiline={true}
                        mode="outlined" // Flat input with an outline
                        style={styles.input}
                        returnKeyType="done"
                        onSubmitEditing={Keyboard.dismiss}
                    />
                    </Card.Content>
                    <Card.Actions style={styles.actions}>
                        <Button style={styles.button} mode="contained" onPress={submitRequest}>Submit</Button>
                    </Card.Actions>
                </Card>
            </View>
        </SafeAreaView>
    </TouchableWithoutFeedback>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 50,
    alignItems: 'center',
  },
  header: {
    fontSize: 20,
    fontWeight: 'bold',
    padding: 16,
    alignItems: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
  },
  backActionPlaceholder: {
    width: 48, 
    height: 48,
  },
  backAction: {
    marginLeft: 0 
  },
  button: {
    textAlign: 'center',
    flex: 1, // Take full width of the card actions
    justifyContent: 'center', // Center the button text
    paddingVertical: 10, // Increase vertical padding
  //   position: 'absolute', 
  //   left: 80,
  //   right: 80, //position, left, right together controls the button's length and horizontal location
  //   alignSelf: 'center', 
  },
  card: {
    marginTop: 16, // Add top margin to push down the card from header
    width: '90%', // Set width to 90% of the screen width
    alignSelf: 'center', // Center the card horizontally
    borderRadius: 8, // Match rounded corners to design
    elevation: 2, // Adjust shadow to match design
  },
  cardTitle: {
    fontSize: 18, // Adjust the font size to make it larger
    fontWeight: 'bold', // Make the text bold
    textAlign: 'center', // Center the text horizontally
    marginTop: 15,
  },
  input: {
    minHeight: 150, // Increased height based on design
    textAlignVertical: 'top', // Align text to the top on Android
    borderRadius: 5, // Reduced rounded corners
  },
  actions: {
    justifyContent: 'center', // Center the button in the actions area
    padding: 16, // Add padding around the button
  },
  request_container: {
    // Container for the card
    width: '100%', // Ensure it takes the full width of the screen
    flex: 1, // Take up remaining space
    justifyContent: 'center', // Center vertically
  },
});

export default PostRequestInfoScreen;
