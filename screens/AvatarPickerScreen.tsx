import React, { useState } from 'react';
import { View, StyleSheet, TouchableOpacity, Image, ScrollView } from 'react-native';
import { Text, Button, Appbar, Card } from 'react-native-paper';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useNavigation } from '@react-navigation/native';
import { useDispatch } from 'react-redux';
import { setAvatarUri } from '../store';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

type RootStackParamList = {
    SettingsScreen: {newAvatarUri: string};
    AvatarPickerScreen: undefined;
  };

type SettingsScreenProps = NativeStackScreenProps<RootStackParamList, 'AvatarPickerScreen'>;

type Avatar = {
    uri: string;
    image: NodeRequire;
}

const avatars: Avatar[] = [
    { uri: '../assets/avatars/avatar1.png', image: require('../assets/avatars/avatar1.png') },
    { uri: '../assets/avatars/avatar2.png', image: require('../assets/avatars/avatar2.png') },
    { uri: '../assets/avatars/avatar3.png', image: require('../assets/avatars/avatar3.png') },
    { uri: '../assets/avatars/avatar4.png', image: require('../assets/avatars/avatar4.png') },
    { uri: '../assets/avatars/avatar5.png', image: require('../assets/avatars/avatar5.png') },
    { uri: '../assets/avatars/avatar6.png', image: require('../assets/avatars/avatar6.png') },
    { uri: '../assets/avatars/avatar7.png', image: require('../assets/avatars/avatar7.png') },
    { uri: '../assets/avatars/avatar8.png', image: require('../assets/avatars/avatar8.png') },
    { uri: '../assets/avatars/avatar9.png', image: require('../assets/avatars/avatar9.png') },
];

let selectedAvatarUri: string = ''

const AvatarPickerScreen: React.FC<SettingsScreenProps> = ({ navigation }: SettingsScreenProps) => {

    const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

    const use_navigation = useNavigation(); //for Appbar.BackAction

    const [selectedAvatar, setSelectedAvatar] = useState<string | null>(null);

    const handleAvatarSelect = (avatar: Avatar) => {
        setSelectedAvatar(avatar.uri);  // Store the URI of the selected avatar
    };

    const dispatch = useDispatch(); 

    const handleSave = (newAvatarUri: string) => {
        if (selectedAvatar) {
            let selectedAvatarUri: string = selectedAvatar // to prevent the "or null" state
            navigation.navigate('SettingsScreen', {newAvatarUri: selectedAvatarUri,}) //useRoute for local navigation
            dispatch(setAvatarUri(selectedAvatarUri)); //Redux for global storage 
        }
    };

    return (
        <ScrollView contentContainerStyle={styles.container}>
            <View style={styles.headerContainer}>
                <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
                <Text style={styles.header}>GiveGetGo</Text>
                <View style={styles.backActionPlaceholder} />
            </View>
            <Card style={styles.card}>
                <Card.Title title="Choose Your Avatar" titleStyle={styles.title} />
                <Card.Content>
                <View style={styles.avatarContainer}>
                    {avatars.map((avatar, index) => (
                    <TouchableOpacity key={index} onPress={() => handleAvatarSelect(avatar)}>
                        <Image
                        source={avatar.image}
                        style={[
                            styles.avatar,
                            selectedAvatar === avatar.uri ? styles.selectedAvatar : {},
                        ]}
                        />
                    </TouchableOpacity>
                    ))}
                </View>
                <Button style={styles.button} mode="contained" onPress={handleSave} disabled={!selectedAvatar}>
                    Save
                </Button>
                </Card.Content>
            </Card>
        </ScrollView>
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
      card: { //page gets longer when there are more contexts
        borderRadius: 15, // Add rounded corners to the card
        marginVertical: 6,
        marginHorizontal: 20,
        elevation: 0, // Adjust for desired shadow depth
        // backgroundColor: '#ffffff', 
        padding: 15, // Add padding inside the card
        // height: 130,
      },
    title: {
        fontSize: 20,
        fontWeight: '500',
        alignSelf: 'center',
        marginTop: 15,
        marginBottom: 15,
    },
    button: {
        textAlign: 'center',
        alignSelf: 'center', 
        marginBottom: -5,
      },
    avatarContainer: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        justifyContent: 'center',
        marginBottom: 20,
    },
    avatar: {
        width: 70,
        height: 70,
        margin: 10,
        opacity: 0.6,
    },
    selectedAvatar: {
        opacity: 1,
        borderWidth: 3, 
        borderColor: '#4CAF50', //light green
        borderRadius: 35, 
    },
});

export default AvatarPickerScreen;