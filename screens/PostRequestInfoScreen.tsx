import React, { useState, useEffect } from 'react';
import { View, StyleSheet, SafeAreaView, Keyboard, TouchableWithoutFeedback  } from 'react-native';
import { Button, Text, Card, TextInput, Appbar } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useRoute } from '@react-navigation/native';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

type RootStackParamList = {
  HomeScreen: undefined;
  PostRequestInfoScreen: undefined;
  PostRequestSucceedScreen: { postId: string, name: string };
};

type PostBid = { 
  user_id: string;
  message: string;
}

const defaultPostBid: PostBid = {
  user_id: '000',
  message: '',
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

const PostRequestInfoScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const route = useRoute();
  const { postId } = route.params as { postId: string };
  const { name } = route.params as { name: string };

  console.log("postID parsed to PostRequestInfoScreen: ", postId) 
  console.log("postOwnerName parsed to PostRequestInfoScreen: ", name) 

  const [postBid, setPostBid] = useState<PostBid>(defaultPostBid);

  useEffect(() => {                                                                                     //fill this in to get db info 
    const fetchPostBid = async () => {
      try {
        const response = await fetch('URL_TO_YOUR_BACKEND/bid/{postId}');                               // postId is used here
        const json = await response.json();
        setPostBid(json); // Adjust this depending on the structure of your JSON
      } catch (error) {
        // console.error(error); // uncomment this after finish frontend developing
      }
    };
    fetchPostBid();
  }, []);

  const use_navigation = useNavigation(); //for Appbar.BackAction

  const [request, setRequest] = React.useState('');

  // update postBid.message to the received "request"
  function updatePostBid(currentInfo: PostBid, updates: Partial<PostBid>): PostBid {
    return {
      ...currentInfo,
      ...updates,
    };
  }

  const submitRequest = (postId: string, name: string) => {
    setPostBid(prev => updatePostBid(prev, { message: request }))
    console.log('Bid sent:', request);
    navigation.navigate('PostRequestSucceedScreen', {postId: postId, name: name});
  };

  useEffect(() => {
    // This will log the updated rating only when postBid.message changes
    console.log('Updated Rating in postOwnerProfile: ', postBid.message);
  }, [postBid.message]);

  // Need to submit the edited postBid (with updated postBid.message) to backend; MAKE SURE the api calling is inside useEffect so that the new value could be fetched

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
                    <Card.Title title="Bid Info" titleStyle={styles.cardTitle}/>
                    <Card.Content>
                    <TextInput
                        label="Please provide details here..."
                        value={request}
                        onChangeText={setRequest}
                        multiline={true}
                        numberOfLines={4}
                        maxLength={150}
                        style={styles.input}
                        returnKeyType="done"
                        onSubmitEditing={Keyboard.dismiss}
                    />
                    </Card.Content>
                    <Card.Actions style={styles.actions}>
                        <Button
                          style={styles.button}
                          mode="contained"
                          onPress={() => submitRequest(postId, name)} 
                        >
                          Submit
                        </Button>
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
    justifyContent: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly horizontally
    paddingLeft: 10, 
    paddingRight: 10, 
    position: 'absolute', // So that while setting card to the vertical middle, it still stays at the same place
    top: 0, 
    left: 0,
    right: 0,
    zIndex: 1, // Ensure the headerContainer is above the card
  },
  header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontFamily: 'Montserrat_700Bold_Italic',
    textAlign: 'center', // Center the text
    color: '#444444', // Dark gray color
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  button: {
    textAlign: 'center',
    flex: 1, // Take full width of the card actions
    justifyContent: 'center', // Center the button text
    paddingVertical: 4, // Increase vertical padding
  },
  card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 15, // Add padding inside the card
    // marginTop: 170,
  },
  cardTitle: {
    fontSize: 18, // Adjust the font size to make it larger
    fontWeight: 'bold', // Make the text bold
    textAlign: 'center', // Center the text horizontally
    marginTop: 6,
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
